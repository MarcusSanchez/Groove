package actions

import (
	. "GrooveGuru/pkgs/util"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// GetArtist returns the artist with the given id.
func (*Actions) GetArtist(c *fiber.Ctx, artistID string) error {
	return proxyArtistRequest(c, Proxy{
		Access:   c.Locals("access").(string),
		Endpoint: "/artists/" + artistID,
	})
}

// GetRelatedArtists returns the artists related to the artist with the given id.
func (*Actions) GetRelatedArtists(c *fiber.Ctx, artistID string) error {
	return proxyArtistRequest(c, Proxy{
		Endpoint: "/artists/" + artistID + "/related-artists",
		Access:   c.Locals("access").(string),
	})
}

// GetArtistTopTracks returns the top tracks of the artist with the given id.
func (*Actions) GetArtistTopTracks(c *fiber.Ctx, artistID string) error {
	return proxyArtistRequest(c, Proxy{
		Endpoint: "/artists/" + artistID + "/top-tracks?market=US",
		Access:   c.Locals("access").(string),
	})
}

// GetArtistAlbums returns the albums of the artist with the given id.
func (*Actions) GetArtistAlbums(c *fiber.Ctx, artistID string) error {
	return proxyArtistRequest(c, Proxy{
		Endpoint: "/artists/" + artistID + "/albums?market=US&limit=50&include_groups=album",
		Access:   c.Locals("access").(string),
	})
}

// proxyArtistRequest proxies a request to the Spotify API for an artist.
// Returns 200 with the artist data if successful.
// Returns 400 if the artist-id is invalid.
// Returns 404 if the artist is not found.
func proxyArtistRequest(c *fiber.Ctx, proxy Proxy) error {
	resp, err := resty.New().R().
		SetHeaders(Headers{
			"Authorization": "Bearer " + proxy.Access,
			"Accept":        "application/json",
		}).
		Get("https://api.spotify.com/v1" + proxy.Endpoint)
	if err != nil {
		LogError("proxyArtistRequest", "Requesting "+proxy.Endpoint, err)
		return InternalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 200:
		c.Set("Content-Type", "application/json")
		return c.Status(fiber.StatusOK).Send(resp.Body())
	case 400:
		return BadRequest(c, "invalid artist-id")
	case 404:
		return BadRequest(c, "artist not found", 404)
	default:
		LogError(
			"proxyArtistRequest",
			"Requesting "+c.Path(),
			errors.New(strconv.Itoa(resp.StatusCode())+": "+string(resp.Body())),
		)
		return InternalServerError(c, "error requesting "+c.Path())
	}
}
