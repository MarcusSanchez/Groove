package actions

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	. "groove/pkgs/util"
	"strconv"
)

// GetTrack returns a track object from the Spotify API.
func (a *Actions) GetTrack(c *fiber.Ctx, trackID string) error {
	return proxyTrackRequest(c, Proxy{
		Endpoint: "/tracks/" + trackID + "?market=US",
		Access:   c.Locals("access").(string),
	})
}

// proxyTrackRequest proxies a request to the Spotify API for track.
// returns 400 if the track-id is invalid.
// returns 404 if the track is not found.
// returns 200 if the track is found.
func proxyTrackRequest(c *fiber.Ctx, proxy Proxy) error {
	resp, err := resty.New().R().
		SetHeaders(Headers{
			"Authorization": "Bearer " + proxy.Access,
			"Accept":        "application/json",
		}).
		Get("https://api.spotify.com/v1" + proxy.Endpoint)
	if err != nil {
		LogError("proxyTrackRequest", "Requesting "+proxy.Endpoint, err)
		return InternalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 200:
		c.Set("Content-Type", "application/json")
		return c.Status(fiber.StatusOK).Send(resp.Body())
	case 400:
		return BadRequest(c, "invalid track-id")
	case 404:
		return BadRequest(c, "track not found", 404)
	default:
		LogError(
			"proxyTrackRequest",
			"Requesting "+c.Path(),
			errors.New(strconv.Itoa(resp.StatusCode())+": "+string(resp.Body())),
		)
		return InternalServerError(c, "error requesting "+c.Path())
	}
}
