package actions

import (
	"GrooveGuru/pkgs/ent"
	OAuthState "GrooveGuru/pkgs/ent/oauthstate"
	SpotifyLink "GrooveGuru/pkgs/ent/spotifylink"
	. "GrooveGuru/pkgs/util"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	accessType = "offline"
	scopes     = []string{
		"playlist-read-private",
		"playlist-read-collaborative",
		"playlist-modify-public",
		"playlist-modify-private",
		"user-library-read",
		"user-library-modify",
	}
)

// LinkSpotify creates a SpotifyLink and sends Spotify
// Authorization page that the client will redirect the user to.
// Returns 200 if successful.
func (a *Actions) LinkSpotify(c *fiber.Ctx) error {
	session := c.Locals("session").(*ent.Session)
	ctx := c.Context()

	// generate a random 16 character string for the state parameter.
	state := strings.ReplaceAll(uuid.New().String(), "-", "")[:16]

	// invalidate any previous state for this user.
	_, err := a.client.OAuthState.
		Delete().
		Where(OAuthState.UserIDEQ(session.UserID)).
		Exec(ctx)
	if err != nil {
		LogError("LinkSpotify", "Checking state", err)
		return InternalServerError(c, "error linking spotify")
	}

	// set state in OAuth-Store for later verification.
	_, err = a.client.OAuthState.Create().
		SetState(state).
		SetExpiration(time.Now().Add(30 * time.Minute)).
		SetUserID(session.UserID).
		Save(ctx)
	if err != nil {
		LogError("LinkSpotify", "Creating OAuthState", err)
		return InternalServerError(c, "error linking spotify")
	}

	baseURL, _ := url.Parse("https://accounts.spotify.com/authorize")
	baseURL.RawQuery = URLSearchParams(Params{
		"response_type": "code",
		"client_id":     a.env.SpotifyClient,
		"scopes":        strings.Join(scopes, " "),
		"redirect_uri":  a.env.BackendURL + "/api/spotify/callback",
		"state":         state,
		"access_type":   accessType,
	})

	return c.Status(200).SendString(baseURL.String())
}

// SpotifyCallback handles the redirect from the Spotify Authorization page.
// It verifies the state and retrieves the access token and refresh token using the code.
// It then saves the tokens as a SpotifyLink, successfully linking Groove and Spotify accounts.
// Returns 201 if successful.
func (a *Actions) SpotifyCallback(c *fiber.Ctx, code, state string) error {
	session := c.Locals("session").(*ent.Session)
	ctx := c.Context()

	// verify state to prevent CSRF.
	storedState, err := a.client.OAuthState.
		Query().
		Where(OAuthState.UserIDEQ(session.UserID)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return Unauthorized(c, "unidentified state")
		}
		LogError("SpotifyCallback", "Checking state", err)
		return InternalServerError(c, "error linking spotify")
	}

	if storedState.Expiration.Before(time.Now()) {
		return Unauthorized(c, "expired state")
	}

	// check if state matches.
	if storedState.State != state {
		LogError(
			"SpotifyCallback",
			"Potential CSRF Attempt",
			errors.New("state mismatch for user: "+strconv.Itoa(session.UserID)),
		)
		return Forbiddened(c, "state mismatch")
	}

	// clear state from store
	err = a.client.OAuthState.DeleteOne(storedState).Exec(ctx)
	if err != nil {
		LogError("SpotifyCallback", "Deleting state", err)
		// no need to alert client background worker will handle it.
	}

	// create base64 encoded string of client id and secret. (as per spotify docs)
	credentials := a.env.SpotifyClient + ":" + a.env.SpotifySecret
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))

	// retrieve access token and refresh token from spotify.
	resp, err := resty.New().R().
		SetHeaders(Headers{
			"Content-Type":  "application/x-www-Form-urlencoded",
			"Authorization": "Basic " + encodedCredentials,
		}).
		SetFormData(Form{
			"code":         code,
			"redirect_uri": a.env.BackendURL + "/api/spotify/callback",
			"grant_type":   "authorization_code",
		}).
		Post("https://accounts.spotify.com/api/token")
	if err != nil {
		LogError("SpotifyCallback", "Requesting token", err)
		return InternalServerError(c, "error requesting token")
	}

	type TokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	payload := new(TokenResponse)
	if json.Unmarshal(resp.Body(), payload) != nil {
		LogError("SpotifyCallback", "Unmarshalling token", err)
		return InternalServerError(c, "error unmarshalling token")
	}

	// save access token and refresh token as SpotifyLink.
	_, err = a.client.SpotifyLink.Create().
		SetAccessToken(payload.AccessToken).
		// Spotify's Access-Token expire after 1 hour, so we set the expiration to 58 minutes to be safe.
		SetAccessTokenExpiration(time.Now().Add(Time58Minutes)).
		SetRefreshToken(payload.RefreshToken).
		SetUserID(session.UserID).
		Save(ctx)
	if err != nil {
		LogError("SpotifyCallback", "Creating spotify link", err)
		return InternalServerError(c, "error linking spotify")
	}

	return c.Redirect(a.env.FrontendURL+"/dashboard/profile", 302)
}

// UnlinkSpotify deletes the SpotifyLink for the user.
// Returns 204 no content if successful.
func (a *Actions) UnlinkSpotify(c *fiber.Ctx) error {
	session := c.Locals("session").(*ent.Session)
	ctx := c.Context()

	// ensure user has a linked spotify account.
	link, err := a.client.SpotifyLink.
		Query().
		Where(SpotifyLink.UserIDEQ(session.UserID)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return Unauthorized(c, "account not linked")
		}
		LogError("UnlinkSpotify", "Checking spotify link", err)
		return InternalServerError(c, "error unlinking spotify")
	}

	// delete spotify link from database.
	if err = a.client.SpotifyLink.DeleteOne(link).Exec(ctx); err != nil {
		LogError("UnlinkSpotify", "Deleting spotify link", err)
		return InternalServerError(c, "error unlinking spotify")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetCurrentUser retrieves the current user's ID if they are linked to Spotify.
// required to grab the user's playlists.
// Returns 200 if successful.
func (a *Actions) GetCurrentUser(c *fiber.Ctx) error {
	access := c.Locals("access").(string)

	// grab current user's ID from Spotify.
	resp, err := resty.New().R().
		SetHeaders(Headers{
			"Authorization": "Bearer " + access,
		}).
		Get("https://api.spotify.com/v1/me")
	if err != nil {
		LogError("GetCurrentUser", "Requesting current user", err)
		return InternalServerError(c, "error getting current user")
	}

	type CurrentUserResponse struct {
		ID string `json:"id"`
	}

	payload := new(CurrentUserResponse)
	if json.Unmarshal(resp.Body(), payload) != nil {
		LogError("GetCurrentUser", "Unmarshalling current user", err)
		return InternalServerError(c, "error getting current user")
	}

	return c.Status(fiber.StatusOK).JSON(payload)
}
