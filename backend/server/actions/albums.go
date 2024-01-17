package actions

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	. "groove/pkgs/util"
	"strconv"
)

// GetAlbum returns the album with the given id.
func (*Actions) GetAlbum(c *fiber.Ctx, albumID string) error {
	return proxyAlbumRequest(c, Proxy{
		Endpoint: "/albums/" + albumID,
		Access:   c.Locals("access").(string),
	})
}

// GetAlbumTracks returns the tracks of the album with the given id.
func (*Actions) GetAlbumTracks(c *fiber.Ctx, albumID string) error {
	return proxyAlbumRequest(c, Proxy{
		Endpoint: "/albums/" + albumID + "/tracks?limit=50&market=US",
		Access:   c.Locals("access").(string),
	})
}

// proxyAlbumRequest proxies a request to the Spotify API for an album.
// returns 200 with the album tracks if successful.
// returns 400 if the album-id is invalid.
// returns 404 if the album is not found.
func proxyAlbumRequest(c *fiber.Ctx, proxy Proxy) error {
	resp, err := resty.New().R().
		SetHeaders(Headers{
			"Authorization": "Bearer " + proxy.Access,
			"Accept":        "application/json",
		}).
		Get("https://api.spotify.com/v1" + proxy.Endpoint)
	if err != nil {
		LogError("proxyAlbumRequest", "Requesting "+proxy.Endpoint, err)
		return InternalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 200:
		c.Set("Content-Type", "application/json")
		return c.Status(fiber.StatusOK).Send(resp.Body())
	case 400:
		return BadRequest(c, "invalid album-id")
	case 404:
		return BadRequest(c, "album not found", 404)
	default:
		LogError(
			"proxyAlbumRequest",
			"Requesting "+c.Path(),
			errors.New(strconv.Itoa(resp.StatusCode())+": "+string(resp.Body())),
		)
		return InternalServerError(c, "error requesting "+c.Path())
	}
}
