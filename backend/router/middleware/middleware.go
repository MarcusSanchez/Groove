package middleware

import (
	"GrooveGuru/db"
	"GrooveGuru/ent"
	Session "GrooveGuru/ent/session"
	SpotifyLink "GrooveGuru/ent/spotifylink"
	"GrooveGuru/env"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recovery "github.com/gofiber/fiber/v2/middleware/recover"
	"net/url"
	"strconv"
	"time"
)

var client = db.Instance()

type (
	form    map[string]string
	headers map[string]string
	params  map[string]string
)

// Attach attaches the middleware that run on all endpoints.
func Attach(app *fiber.App) {
	app.Static("/", "./public")
	// catch-all route for the frontend.
	app.Use("/", ReactServer)
	app.Use(logger.New())
	// if the server were to crash, this would restart the server.
	app.Use(recovery.New())
	switch env.IsProd {
	case false:
		// in development, frontend and backend are listening on different ports;
		// therefore CORS needs to be configured to allow the frontend url.
		app.Use(cors.New(cors.Config{
			AllowOrigins:     env.FrontendURL,
			AllowCredentials: true,
		}))
	case true:
		// limits repeated requests to endpoints; protection against brute-force attacks.
		app.Use(limiter.New())
	}
}

// ReactServer serves the frontend.
// this is used for the catch-all route.
// if route starts with /api, it will not be served by this function.
func ReactServer(c *fiber.Ctx) error {
	path := c.Path()
	if len(path) > 4 && path[:4] == "/api" {
		return c.Next()
	}
	return c.SendFile("./public/index.html")
}

// RedirectAuthorized redirects to the home page if the user is authorized.
// Useful for instances where the user should not be logged-in/authenticated.
// i.e. login and register pages.
func RedirectAuthorized(c *fiber.Ctx) error {
	ctx := c.Context()

	authorization := c.Cookies("Authorization")
	if authorization == "" {
		return c.Next()
	}

	// check if cookie session actually exists.
	session, err := client.Session.
		Query().
		Where(Session.TokenEQ(authorization)).
		First(ctx)
	if ent.IsNotFound(err) {
		expireSessionCookies(c)
		return c.Next()
	} else if err != nil {
		logError("RedirectAuthorized[MIDDLEWARE]", "checking session", err)
		return c.Status(fiber.StatusInternalServerError).SendString("error while authorizing")
	}

	// check if session has expired.
	if session.Expiration.Before(time.Now()) {
		expireSessionCookies(c)
		return c.Next()
	}

	return c.SendStatus(fiber.StatusPermanentRedirect)
}

// AuthorizeAny authorizes the user if the Authorization cookie is set and valid (no permissions necessary).
// for general use of endpoints where the user just needs to be logged in.
// i.e. viewing content, searching songs, etc...
func AuthorizeAny(c *fiber.Ctx) error {
	ctx := c.Context()

	authorization := c.Cookies("Authorization")
	if authorization == "" {
		return unauthorized(c)
	}

	// check if cookie session actually exists.
	session, err := client.Session.
		Query().
		Where(Session.TokenEQ(authorization)).
		First(ctx)
	if ent.IsNotFound(err) {
		expireSessionCookies(c)
		return unauthorized(c)
	} else if err != nil {
		logError("AuthorizeAny[MIDDLEWARE]", "checking session", err)
		return c.Status(fiber.StatusInternalServerError).SendString("error while authorizing")
	}

	// check if session has expired.
	if session.Expiration.Before(time.Now()) {
		expireSessionCookies(c)
		return unauthorized(c)
	}

	c.Locals("session", session)
	return c.Next()
}

// RedirectLinked redirects to the home page if the user is already linked to spotify.
// Useful for instances where the user should not be linked to spotify.
// i.e. spotify link page.
//
// NOTE: this middleware is meant to be used after with AuthorizeAny or CheckCSRF to retrieve session.
func RedirectLinked(c *fiber.Ctx) error {
	session := c.Locals("session").(*ent.Session)
	ctx := c.Context()

	// check if a link already exists.
	exists, err := client.SpotifyLink.
		Query().
		Where(SpotifyLink.UserIDEQ(session.UserID)).
		Exist(ctx)
	if err != nil {
		logError("RedirectLinked[MIDDLEWARE]", "checking spotify link", err)
		return c.Status(fiber.StatusInternalServerError).SendString("error while authorizing")
	} else if !exists {
		return c.Next()
	}

	return c.SendStatus(fiber.StatusPermanentRedirect)
}

// CheckCSRF checks if the Csrf token in the body matches the one in the session.
// for use of endpoints where the user needs to be logged in and the request needs to be verified.
// i.e. actions with side effects (creating, updating, deleting, etc...)
//
// NOTE: this middleware fulfills the same purpose as AuthorizeAny, no need to use both.
func CheckCSRF(c *fiber.Ctx) error {
	ctx := c.Context()

	type CSRF struct {
		Csrf string `json:"csrf_"`
	}

	authorization := c.Cookies("Authorization")
	if authorization == "" {
		return unauthorized(c)
	}

	var payload CSRF
	if c.BodyParser(&payload) != nil {
		return forbiddened(c)
	}

	// check if cookie session actually exists.
	session, err := client.Session.
		Query().
		Where(Session.TokenEQ(authorization)).
		First(ctx)
	if ent.IsNotFound(err) {
		expireSessionCookies(c)
		return unauthorized(c)
	} else if err != nil {
		logError("CheckCSRF[MIDDLEWARE]", "checking session", err)
		return c.Status(fiber.StatusInternalServerError).SendString("error while authorizing")
	}

	// check if session has expired.
	if session.Expiration.Before(time.Now()) {
		expireSessionCookies(c)
		return unauthorized(c)
	}

	if session.Csrf != payload.Csrf {
		// request was forged.
		logError(
			"CheckCSRF[MIDDLEWARE]",
			"Potential CSRF Attempt",
			errors.New("csrf token mismatch for user: "+strconv.Itoa(session.UserID)),
		)
		return forbiddened(c)
	}

	c.Locals("session", session)
	return c.Next()
}

// SetAccess sets the access token for the spotify client.
// if user is linked to spotify, the access token will be theirs;
// otherwise, the access token will be the default token.
func SetAccess(c *fiber.Ctx) error {
	session := c.Locals("session").(*ent.Session)
	ctx := c.Context()

	// check if the user is linked to spotify, if so, use their access token.
	link, err := client.SpotifyLink.
		Query().
		Where(SpotifyLink.UserIDEQ(session.UserID)).
		First(ctx)
	if ent.IsNotFound(err) {
		// if not, use the default access token.
		link, err = defaultAccessToken()
		if err != nil {
			logError("SetAccess[MIDDLEWARE]", "getting default access token", err)
			return c.Status(fiber.StatusInternalServerError).SendString("error while authorizing")
		}
	} else if err != nil {
		logError("SetAccess[MIDDLEWARE]", "checking spotify link", err)
		return c.Status(fiber.StatusInternalServerError).SendString("error while authorizing")
	}

	// if the link hasn't expired, we can use it.
	if !link.AccessTokenExpiration.Before(time.Now()) {
		c.Locals("access", link.AccessToken)
		return c.Next()
	}

	// otherwise, we need to refresh the token.
	credentials := env.SpotifyClient + ":" + env.SpotifySecret
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))

	http := resty.New()
	resp, err := http.R().
		SetHeaders(headers{
			"Content-Type":  "application/x-www-form-urlencoded",
			"Authorization": "Basic " + encodedCredentials,
		}).
		SetBody(urlSearchParams(params{
			"grant_type":    "refresh_token",
			"refresh_token": link.RefreshToken,
		})).
		Post("https://accounts.spotify.com/api/token")
	if err != nil {
		logError("SetAccess[MIDDLEWARE]", "refreshing token", err)
		return c.Status(fiber.StatusInternalServerError).SendString("error while authorizing")
	}

	if resp.StatusCode() != 200 {
		logError(
			"SetAccess[MIDDLEWARE]",
			"refreshing token",
			errors.New(fmt.Sprintln(resp.StatusCode(), ", ", string(resp.Body()))),
		)
		return c.Status(fiber.StatusInternalServerError).SendString("error while authorizing")
	}

	type TokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	payload := TokenResponse{}
	if json.Unmarshal(resp.Body(), &payload) != nil {
		logError("SetAccess[MIDDLEWARE]", "unmarshalling token", err)
		return c.Status(fiber.StatusInternalServerError).SendString("error while authorizing")
	}

	if payload.RefreshToken == "" {
		payload.RefreshToken = link.RefreshToken
	}

	// update link with new access and refresh tokens.
	_, err = client.SpotifyLink.UpdateOne(link).
		SetAccessToken(payload.AccessToken).
		SetRefreshToken(payload.RefreshToken).
		SetAccessTokenExpiration(time.Now().Add(58 * time.Minute)).
		Save(ctx)
	if ent.IsValidationError(err) {
		logError("SetAccess[MIDDLEWARE]", "updating spotify link (validation)", err)
		return c.Status(fiber.StatusInternalServerError).SendString("error while authorizing")
	} else if err != nil {
		logError("SetAccess[MIDDLEWARE]", "updating spotify link", err)
		return c.Status(fiber.StatusInternalServerError).SendString("error while authorizing")
	}

	c.Locals("access", payload.AccessToken)
	return c.Next()
}

// AuthorizeLinked authorizes the user if the Authorization cookie is set and valid (no permissions necessary).
// for general use of endpoints where the user just needs to be logged in and linked with spotify.
// i.e. viewing, editing playlists, etc...
func AuthorizeLinked(c *fiber.Ctx) error {
	ctx := c.Context()

	authorization := c.Cookies("Authorization")
	if authorization == "" {
		return unauthorized(c)
	}

	// check if cookie session actually exists.
	session, ok := c.Locals("session").(*ent.Session)
	if !ok {
		var err error
		session, err = client.Session.Query().Where(Session.TokenEQ(authorization)).First(ctx)
		if ent.IsNotFound(err) {
			expireSessionCookies(c)
			return unauthorized(c)
		} else if err != nil {
			logError("AuthorizeAny[MIDDLEWARE]", "checking session", err)
			return c.Status(fiber.StatusInternalServerError).SendString("error while authorizing")
		}

		// check if session has expired.
		if session.Expiration.Before(time.Now()) {
			expireSessionCookies(c)
			return unauthorized(c)
		}
	}

	// check if user is linked, else reject the request
	exists, err := client.SpotifyLink.
		Query().
		Where(SpotifyLink.UserIDEQ(session.UserID)).
		Exist(ctx)
	if err != nil {
		logError("AuthorizeLinked[MIDDLEWARE]", "checking spotify link", err)
		return c.Status(fiber.StatusInternalServerError).SendString("error while authorizing")
	} else if !exists {
		return forbiddened(c)
	}

	c.Locals("session", session)
	return c.Next()
}

/** Helpers **/

func forbiddened(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).SendString("Forbidden")
}

func unauthorized(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
}

func logError(fn, context string, err error) {
	fmt.Printf(
		"%s [ERROR] [Function: %s (Context: %s)] %s\n",
		time.Now().Format("15:04:05"),
		fn, context, err.Error(),
	)
}

// expireSessionCookies deletes the Authorization and Csrf cookies.
//
// This is used over ClearCookie because:
// Web browsers and other compliant clients will only clear the cookie
// if the given options are identical to those when creating the cookie.
func expireSessionCookies(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		SameSite: env.SameSite,
		Secure:   env.Secure,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "Csrf",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: false,
		SameSite: env.SameSite,
		Secure:   env.Secure,
	})
}

func urlSearchParams(params map[string]string) string {
	qParams := url.Values{}
	for key, value := range params {
		qParams.Add(key, value)
	}
	return qParams.Encode()
}

func defaultAccessToken() (*ent.SpotifyLink, error) {
	link, err := client.SpotifyLink.
		Query().
		Where(SpotifyLink.UserIDEQ(1)).
		First(context.Background())
	if ent.IsNotFound(err) {
		panic("default access token not found")
	} else if err != nil {
		return nil, err
	}

	return link, err
}
