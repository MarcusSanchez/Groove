package middleware

import (
	"GrooveGuru/db"
	"GrooveGuru/ent"
	Session "GrooveGuru/ent/session"
	"GrooveGuru/env"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recovery "github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"time"
)

var client, ctx = db.Instance()

// Attach attaches the middleware that run on all endpoints
func Attach(app *fiber.App) {
	app.Use(logger.New())
	app.Use(recovery.New())
	if !env.IsProd {
		// in development, frontend and backend are listening on different ports;
		// therefore CORS needs to be configured to allow all origins.
		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
		}))
	}
}

// RedirectAuthorized redirects to the home page if the user is authorized.
func RedirectAuthorized(c *fiber.Ctx) error {
	authorization := c.Cookies("Authorization")
	if authorization == "" {
		return c.Next()
	}

	// check if cookie session actually exists.
	session, err := client.Session.Query().Where(Session.TokenEQ(authorization)).First(ctx)
	if ent.IsNotFound(err) {
		c.ClearCookie("Authorization")
		c.ClearCookie("Csrf")
		return c.Next()
	} else if err != nil {
		logError("RedirectAuthorized[MIDDLEWARE]", "checking session", err)
		return c.Status(fiber.StatusInternalServerError).SendString("error while authorizing")
	}

	// check if session has expired.
	if session.Expiration.Before(time.Now()) {
		c.ClearCookie("Authorization")
		c.ClearCookie("Csrf")
		return c.Next()
	}

	return c.Status(fiber.StatusPermanentRedirect).Redirect(c.Hostname() + "/")
}

// AuthorizeAny authorizes the user if the Authorization cookie is set and valid (no permissions necessary).
func AuthorizeAny(c *fiber.Ctx) error {
	authorization := c.Cookies("Authorization")
	if authorization == "" {
		return unauthorized(c)
	}

	// check if cookie session actually exists.
	session, err := client.Session.Query().Where(Session.TokenEQ(authorization)).First(ctx)
	if ent.IsNotFound(err) {
		c.ClearCookie("Authorization")
		c.ClearCookie("Csrf")
		return unauthorized(c)
	} else if err != nil {
		logError("AuthorizeAny[MIDDLEWARE]", "checking session", err)
		return c.Status(fiber.StatusInternalServerError).SendString("error while authorizing")
	}

	// check if session has expired.
	if session.Expiration.Before(time.Now()) {
		c.ClearCookie("Authorization")
		c.ClearCookie("Csrf")
		return unauthorized(c)
	}

	c.Locals("session", session)
	return c.Next()
}

// CheckCSRF checks if the Csrf token in the body matches the one in the session.
func CheckCSRF(c *fiber.Ctx) error {

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
	session, err := client.Session.Query().Where(Session.TokenEQ(authorization)).First(ctx)
	if ent.IsNotFound(err) {
		c.ClearCookie("Authorization")
		c.ClearCookie("Csrf")
		return unauthorized(c)
	} else if err != nil {
		return forbiddened(c)
	}

	if session.Csrf != payload.Csrf {
		//Csrf was forged.
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
	log.Printf(
		"[%s] [ERROR] [Function: %s (Context: %s)] %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		fn, context, err.Error(),
	)
}
