package actions

import (
	"GrooveGuru/db"
	"GrooveGuru/env"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

const week = time.Hour * 24 * 7

var client = db.Instance()

// setSessionCookies sets the Authorization and Csrf cookies.
func setSessionCookies(c *fiber.Ctx, authorization, csrf string, expiration time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    authorization,
		Expires:  expiration,
		HTTPOnly: true,
		SameSite: env.SameSite,
		Secure:   env.Secure,
	})

	c.Cookie(&fiber.Cookie{
		Name:    "Csrf",
		Value:   csrf,
		Expires: expiration,
		// HttpOnly is set to false because the frontend needs to access it.
		// This isn't a security risk because the cookie is for CSRF protection;
		// If XSS is present, the attacker can already do anything they want.
		HTTPOnly: false,
		SameSite: env.SameSite,
		Secure:   env.Secure,
	})
}

// expireSessionCookies expires the Authorization and Csrf cookies.
//
// This is used over ClearCookie because:
// Web browsers and other compliant clients will only clear the cookie
// if the given options are identical to those when creating the cookie
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

// logError formats and prints error with context.
func logError(fn, context string, err error) {
	fmt.Printf(
		"%s [ERROR] [Function: %s (Context: %s)] %s\n",
		time.Now().Format("15:04:05"),
		fn, context, err.Error(),
	)
}

/** responses **/
func badRequest(c *fiber.Ctx, msg string, code ...int) error {
	status := fiber.StatusBadRequest
	if len(code) > 0 {
		status = code[0]
	}
	return c.Status(status).JSON(fiber.Map{
		"error":   "bad request",
		"message": msg,
	})
}

func internalServerError(c *fiber.Ctx, msg string, code ...int) error {
	status := fiber.StatusInternalServerError
	if len(code) > 0 {
		status = code[0]
	}
	return c.Status(status).JSON(fiber.Map{
		"error":   "internal server error",
		"message": msg,
	})
}

func unauthorized(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error":   "unauthorized",
		"message": msg,
	})
}

func forbidden(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error":   "forbidden",
		"message": msg,
	})
}
