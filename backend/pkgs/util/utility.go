package util

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"time"
)

const (
	// TimeWeek is used for session cookies expiration
	TimeWeek = 24 * 7 * time.Hour
	// Time58Minutes is used for Spotify Access Token expiration
	Time58Minutes = 58 * time.Minute
)

const (
	SpotifyAPI         = "https://api.spotify.com/v1"
	SpotifyAccountsAPI = "https://accounts.spotify.com/api"
)

type (
	Form    map[string]string
	Headers map[string]string
	Params  map[string]string
)

// LogError formats and prints error with context.
func LogError(fn, context string, err error) {
	fmt.Printf(
		"%s [ERROR] [Function: %s (Context: %s)] %s\n",
		time.Now().Format("15:04:05"),
		fn, context, err.Error(),
	)
}

// SetSessionCookies sets the Authorization and Csrf cookies.
func SetSessionCookies(c *fiber.Ctx, authorization, csrf string, expiration time.Time, sameSite string, secure bool) {
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    authorization,
		Expires:  expiration,
		HTTPOnly: true,
		SameSite: sameSite,
		Secure:   secure,
	})

	c.Cookie(&fiber.Cookie{
		Name:    "Csrf",
		Value:   csrf,
		Expires: expiration,
		// HttpOnly is set to false because the frontend needs to access it.
		// This isn't a security risk because the cookie is for CSRF protection;
		// If XSS is present, the attacker can already do anything they want.
		HTTPOnly: false,
		SameSite: sameSite,
		Secure:   secure,
	})
}

// ExpireSessionCookies deletes the Authorization and Csrf cookies.
func ExpireSessionCookies(c *fiber.Ctx, sameSite string, secure bool) {
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		SameSite: sameSite,
		Secure:   secure,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "Csrf",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: false,
		SameSite: sameSite,
		Secure:   secure,
	})
}

// URLSearchParams converts a map of query parameters to a URL encoded string.
func URLSearchParams(params map[string]string) string {
	qParams := url.Values{}
	for key, value := range params {
		qParams.Add(key, value)
	}
	return qParams.Encode()
}
