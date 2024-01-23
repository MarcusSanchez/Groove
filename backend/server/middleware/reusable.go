package middleware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/MarcusSanchez/go-parse"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"groove/pkgs/ent"
	Session "groove/pkgs/ent/session"
	SpotifyLink "groove/pkgs/ent/spotifylink"
	. "groove/pkgs/util"
	"net/http"
	"strconv"
	"time"
)

// RedirectAuthorized redirects to the home page if the user is authorized.
// Useful for instances where the user should not be logged-in/authenticated.
// i.e. login and register pages.
func (m *Middlewares) RedirectAuthorized(c *fiber.Ctx) error {
	authorization := c.Cookies("Authorization")
	if authorization == "" {
		return c.Next()
	}

	ctx := c.Context()
	// check if cookie session actually exists.
	session, err := m.client.Session.
		Query().
		Where(Session.TokenEQ(authorization)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			ExpireSessionCookies(c, m.env.SameSite, m.env.Secure)
			return c.Next()
		}
		LogError("RedirectAuthorized[MIDDLEWARE]", "checking session", err)
		return InternalServerError(c, "error while authorizing")
	}

	// check if session has expired.
	if session.Expiration.Before(time.Now()) {
		ExpireSessionCookies(c, m.env.SameSite, m.env.Secure)
		return c.Next()
	}

	return c.SendStatus(http.StatusPermanentRedirect)
}

// AuthorizeAny authorizes the user if the Authorization cookie is set and valid (no permissions necessary).
// for general use of endpoints where the user just needs to be logged in.
// i.e. viewing content, searching songs, etc...
func (m *Middlewares) AuthorizeAny(c *fiber.Ctx) error {
	authorization := c.Cookies("Authorization")
	if authorization == "" {
		return Unauthorized(c, "missing authorization")
	}

	ctx := c.Context()
	// check if cookie session actually exists.
	session, err := m.client.Session.
		Query().
		Where(Session.TokenEQ(authorization)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			ExpireSessionCookies(c, m.env.SameSite, m.env.Secure)
			return Unauthorized(c, "session not found")
		}
		LogError("AuthorizeAny[MIDDLEWARE]", "checking session", err)
		return InternalServerError(c, "error while authorizing")
	}

	// check if session has expired.
	if session.Expiration.Before(time.Now()) {
		ExpireSessionCookies(c, m.env.SameSite, m.env.Secure)
		return Unauthorized(c, "session expired")
	}

	c.Locals("session", session)
	return c.Next()
}

// RedirectLinked redirects to the home page if the user is already linked to spotify.
// Useful for instances where the user should not be linked to spotify.
// i.e. spotify link page.
//
// NOTE: this middleware is meant to be used after with AuthorizeAny or CheckCSRF to retrieve session.
func (m *Middlewares) RedirectLinked(c *fiber.Ctx) error {
	session := c.Locals("session").(*ent.Session)
	ctx := c.Context()

	// check if a link already exists.
	exists, err := m.client.SpotifyLink.
		Query().
		Where(SpotifyLink.UserIDEQ(session.UserID)).
		Exist(ctx)
	if err != nil {
		LogError("RedirectLinked[MIDDLEWARE]", "checking spotify link", err)
		return InternalServerError(c, "error while authorizing")
	} else if !exists {
		return c.Next()
	}

	return c.SendStatus(http.StatusPermanentRedirect)
}

// CheckCSRF checks if the Csrf token in the body matches the one in the session.
// for use of endpoints where the user needs to be logged in and the request needs to be verified.
// i.e. actions with side effects (creating, updating, deleting, etc...)
//
// NOTE: this middleware fulfills the same purpose as AuthorizeAny, no need to use both.
func (m *Middlewares) CheckCSRF(c *fiber.Ctx) error {
	authorization := c.Cookies("Authorization")
	if authorization == "" {
		return Unauthorized(c, "missing authorization")
	}

	ctx := c.Context()
	// check if cookie session actually exists.
	session, err := m.client.Session.
		Query().
		Where(Session.TokenEQ(authorization)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			ExpireSessionCookies(c, m.env.SameSite, m.env.Secure)
			return Unauthorized(c, "session not found")
		}
		LogError("CheckCSRF[MIDDLEWARE]", "checking session", err)
		return InternalServerError(c, "error while authorizing")
	}

	// check if session has expired.
	if session.Expiration.Before(time.Now()) {
		ExpireSessionCookies(c, m.env.SameSite, m.env.Secure)
		return Unauthorized(c, "session expired")
	}

	type CSRF struct {
		Csrf string `json:"csrf_"`
	}

	payload, err := parse.JSON[CSRF](c.Body())
	if err != nil {
		return BadRequest(c, err.Error())
	}

	if session.Csrf != payload.Csrf {
		// request was forged.
		LogError(
			"CheckCSRF[MIDDLEWARE]",
			"Potential CSRF Attempt",
			errors.New("csrf token mismatch for user: "+strconv.Itoa(session.UserID)),
		)
		return Forbidden(c, "csrf token mismatch")
	}

	c.Locals("session", session)
	return c.Next()
}

// SetAccess sets the access token for the spotify client.
// if user is linked to spotify, the access token will be theirs;
// otherwise, the access token will be the default token.
func (m *Middlewares) SetAccess(c *fiber.Ctx) error {
	session := c.Locals("session").(*ent.Session)
	ctx := c.Context()

	// check if the user is linked to spotify, if so, use their access token.
	link, err := m.client.SpotifyLink.
		Query().
		Where(SpotifyLink.UserIDEQ(session.UserID)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			// if not, use the default access token.
			link, err = m.defaultAccessToken()
			if err != nil {
				LogError("SetAccess[MIDDLEWARE]", "getting default access token", err)
				return InternalServerError(c, "error while authorizing")
			}
		} else {
			LogError("SetAccess[MIDDLEWARE]", "checking spotify link", err)
			return InternalServerError(c, "error while authorizing")
		}
	}

	// if the link hasn't expired, we can use it.
	if !link.AccessTokenExpiration.Before(time.Now()) {
		c.Locals("access", link.AccessToken)
		return c.Next()
	}

	// otherwise, we need to refresh the token.
	credentials := m.env.SpotifyClient + ":" + m.env.SpotifySecret
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))

	resp, err := resty.New().R().
		SetHeaders(Headers{
			"Content-Type":  "application/x-www-form-urlencoded",
			"Authorization": "Basic " + encodedCredentials,
		}).
		SetBody(URLSearchParams(Params{
			"grant_type":    "refresh_token",
			"refresh_token": link.RefreshToken,
		})).
		Post(SpotifyAccountsAPI + "/token")
	if err != nil {
		LogError("SetAccess[MIDDLEWARE]", "refreshing token", err)
		return InternalServerError(c, "error while authorizing")
	}

	if resp.StatusCode() != 200 {
		LogError(
			"SetAccess[MIDDLEWARE]",
			"refreshing token",
			errors.New(strconv.Itoa(resp.StatusCode())+": "+string(resp.Body())),
		)
		return InternalServerError(c, "error while authorizing")
	}

	type Tokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	payload := new(Tokens)
	if json.Unmarshal(resp.Body(), payload) != nil {
		LogError("SetAccess[MIDDLEWARE]", "unmarshalling token", err)
		return InternalServerError(c, "error while authorizing")
	}

	if payload.RefreshToken == "" {
		payload.RefreshToken = link.RefreshToken
	}

	// update link with new access and refresh tokens.
	_, err = m.client.SpotifyLink.UpdateOne(link).
		SetAccessToken(payload.AccessToken).
		SetRefreshToken(payload.RefreshToken).
		// Spotify's Access-Token expire after 1 hour, so we set the expiration to 58 minutes to be safe.
		SetAccessTokenExpiration(time.Now().Add(Time58Minutes)).
		Save(ctx)
	if err != nil {
		if ent.IsValidationError(err) {
			LogError("SetAccess[MIDDLEWARE]", "updating spotify link (validation)", err)
			return InternalServerError(c, "error while authorizing")
		}
		LogError("SetAccess[MIDDLEWARE]", "updating spotify link", err)
		return InternalServerError(c, "error while authorizing")
	}

	c.Locals("access", payload.AccessToken)
	return c.Next()
}

// AuthorizeLinked authorizes the user if the Authorization cookie is set and valid (no permissions necessary).
// for general use of endpoints where the user just needs to be logged in and linked with spotify.
// i.e. viewing, editing playlists, etc...
func (m *Middlewares) AuthorizeLinked(c *fiber.Ctx) error {
	ctx := c.Context()

	authorization := c.Cookies("Authorization")
	if authorization == "" {
		return Unauthorized(c, "missing authorization")
	}

	// check if cookie session actually exists.
	session, ok := c.Locals("session").(*ent.Session)
	if !ok {
		var err error
		session, err = m.client.Session.Query().Where(Session.TokenEQ(authorization)).First(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				ExpireSessionCookies(c, m.env.SameSite, m.env.Secure)
				return Unauthorized(c, "session not found")
			}
			LogError("AuthorizeAny[MIDDLEWARE]", "checking session", err)
			return InternalServerError(c, "error while authorizing")
		}

		// check if session has expired.
		if session.Expiration.Before(time.Now()) {
			ExpireSessionCookies(c, m.env.SameSite, m.env.Secure)
			return Unauthorized(c, "session expired")
		}
	}

	// check if user is linked, else reject the request
	exists, err := m.client.SpotifyLink.
		Query().
		Where(SpotifyLink.UserIDEQ(session.UserID)).
		Exist(ctx)
	if err != nil {
		LogError("AuthorizeLinked[MIDDLEWARE]", "checking spotify link", err)
		return InternalServerError(c, "error while authorizing")
	} else if !exists {
		return Forbidden(c, "account not linked")
	}

	c.Locals("session", session)
	return c.Next()
}
