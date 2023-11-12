package actions

import (
	"GrooveGuru/ent"
	OAuthState "GrooveGuru/ent/oauthstate"
	"GrooveGuru/env"
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
	scope      = []string{
		"playlist-read-private",
		"playlist-read-collaborative",
		"playlist-modify-public",
		"playlist-modify-private",
		"user-library-read",
		"user-library-modify",
	}
)

type (
	param   []string
	body    map[string]any
	form    map[string]string
	headers map[string]string
)

// LinkSpotify creates a SpotifyLink and sends Spotify
// Authorization page that the client will redirect the user to.
func LinkSpotify(c *fiber.Ctx) error {
	session := c.Locals("session").(*ent.Session)

	// generate a random 16 character string for the state parameter.
	state := strings.ReplaceAll(uuid.New().String(), "-", "")[:16]

	// invalidate any previous state for this user.
	_, err := client.OAuthState.
		Delete().
		Where(OAuthState.UserIDEQ(session.UserID)).
		Exec(ctx)
	if err != nil {
		logError("LinkSpotify", "Checking state", err)
		return internalServerError(c, "error linking spotify")
	}

	// set state in OAuth-Store for later verification.
	_, err = client.OAuthState.Create().
		SetState(state).
		SetExpiration(time.Now().Add(30 * time.Minute)).
		SetUserID(session.UserID).
		Save(ctx)
	if err != nil {
		logError("LinkSpotify", "Creating OAuthState", err)
		return internalServerError(c, "error linking spotify")
	}

	baseURL, _ := url.Parse("https://accounts.spotify.com/authorize")
	qParams := url.Values{
		"response_type": param{"code"},
		"client_id":     param{env.SpotifyClient},
		"scope":         param{strings.Join(scope, " ")},
		"redirect_uri":  param{env.BackendURL + "/api/spotify/callback"},
		"state":         param{state},
		"access_type":   param{accessType},
	}
	baseURL.RawQuery = qParams.Encode()

	return c.Status(200).SendString(baseURL.String())
}

func SpotifyCallback(c *fiber.Ctx, code, state string) error {
	session := c.Locals("session").(*ent.Session)

	// verify state to prevent CSRF.
	storedState, err := client.OAuthState.
		Query().
		Where(OAuthState.UserIDEQ(session.UserID)).
		First(ctx)
	if ent.IsNotFound(err) {
		return unauthorized(c, "unidentified state")
	} else if err != nil {
		logError("SpotifyCallback", "Checking state", err)
		return internalServerError(c, "error linking spotify")
	}

	if storedState.Expiration.Before(time.Now()) {
		return unauthorized(c, "expired state")
	}

	// check if state matches.
	if storedState.State != state {
		logError(
			"SpotifyCallback",
			"Potential CSRF Attempt",
			errors.New("state mismatch for user: "+strconv.Itoa(session.UserID)),
		)
		return forbidden(c, "state mismatch")
	}

	// clear state from store
	err = client.OAuthState.DeleteOne(storedState).Exec(ctx)
	if err != nil {
		// no need to alert client. background worker will handle it.
		logError("SpotifyCallback", "Deleting state", err)
	}

	// create base64 encoded string of client id and secret. (as per spotify docs)
	credentials := env.SpotifyClient + ":" + env.SpotifySecret
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))

	// retrieve access token and refresh token from spotify.
	http := resty.New()
	resp, err := http.R().
		SetHeaders(headers{
			"Content-Type":  "application/x-www-form-urlencoded",
			"Authorization": "Basic " + encodedCredentials,
		}).
		SetFormData(form{
			"code":         code,
			"redirect_uri": env.BackendURL + "/api/spotify/callback",
			"grant_type":   "authorization_code",
		}).
		Post("https://accounts.spotify.com/api/token")
	if err != nil {
		logError("SpotifyCallback", "Requesting token", err)
		return c.Status(500).SendString("error requesting token")
	}

	type TokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	payload := TokenResponse{}
	if json.Unmarshal(resp.Body(), &payload) != nil {
		logError("SpotifyCallback", "Unmarshalling token", err)
		return internalServerError(c, "error unmarshalling token")
	}

	// save access token and refresh token as SpotifyLink.
	_, err = client.SpotifyLink.Create().
		SetAccessToken(payload.AccessToken).
		// Spotify's Access-Token expire after 1 hour, so we set the expiration to 58 minutes to be safe.
		SetAccessTokenExpiration(time.Now().Add(58 * time.Minute)).
		SetRefreshToken(payload.RefreshToken).
		SetUserID(session.UserID).
		Save(ctx)
	if err != nil {
		logError("SpotifyCallback", "Creating spotify link", err)
		return internalServerError(c, "error linking spotify")
	}

	return c.Status(201).JSON(fiber.Map{
		"acknowledged": true,
		"message":      "Spotify linked",
	})
}
