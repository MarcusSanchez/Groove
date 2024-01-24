package actions

import (
	"github.com/gofiber/fiber/v2"
	. "groove/pkgs/util"
)

// GetTrack returns a track object from the Spotify API.
func (a *Actions) GetTrack(c *fiber.Ctx, trackID string) error {
	proxy := Proxy{
		Endpoint: "/tracks/" + trackID + "?market=US",
		Access:   c.Locals("access").(string),
	}
	return proxy.TrackRequest(c)
}
