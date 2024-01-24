package util

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type Proxy struct {
	Endpoint string
	Access   string
}

// AlbumRequest proxies a request to the Spotify API for an album.
// returns 200 with the album tracks if successful.
// returns 400 if the album-id is invalid.
// returns 404 if the album is not found.
func (p *Proxy) AlbumRequest(c *fiber.Ctx) error {
	resp, err := resty.New().R().
		SetHeaders(Headers{
			"Authorization": "Bearer " + p.Access,
			"Accept":        "application/json",
		}).
		Get(SpotifyAPI + p.Endpoint)
	if err != nil {
		LogError("Proxy-AlbumRequest", "Requesting "+p.Endpoint, err)
		return InternalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 200:
		c.Set("Content-Type", "application/json")
		return c.Status(fiber.StatusOK).Send(resp.Body())
	case 400:
		return BadRequest(c, "invalid album-id")
	case 404:
		return BadRequest(c, "album not found", fiber.StatusNotFound)
	default:
		LogError(
			"Proxy-AlbumRequest",
			"Requesting "+c.Path(),
			errors.New(strconv.Itoa(resp.StatusCode())+": "+string(resp.Body())),
		)
		return InternalServerError(c, "error requesting "+c.Path())
	}
}

// ArtistRequest proxies a request to the Spotify API for an artist.
// Returns 200 with the artist data if successful.
// Returns 400 if the artist-id is invalid.
// Returns 404 if the artist is not found.
func (p *Proxy) ArtistRequest(c *fiber.Ctx) error {
	resp, err := resty.New().R().
		SetHeaders(Headers{
			"Authorization": "Bearer " + p.Access,
			"Accept":        "application/json",
		}).
		Get(SpotifyAPI + p.Endpoint)
	if err != nil {
		LogError("Proxy-ArtistRequest", "Requesting "+p.Endpoint, err)
		return InternalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 200:
		c.Set("Content-Type", "application/json")
		return c.Status(fiber.StatusOK).Send(resp.Body())
	case 400:
		return BadRequest(c, "invalid artist-id")
	case 404:
		return BadRequest(c, "artist not found", fiber.StatusNotFound)
	default:
		LogError(
			"Proxy-ArtistRequest",
			"Requesting "+c.Path(),
			errors.New(strconv.Itoa(resp.StatusCode())+": "+string(resp.Body())),
		)
		return InternalServerError(c, "error requesting "+c.Path())
	}
}

// PlaylistRequest proxies a request to the Spotify API for a playlist.
// returns 200 with the playlist tracks if successful.
// returns 400 if the playlist-id is invalid.
// returns 404 if the playlist is not found.
func (p *Proxy) PlaylistRequest(c *fiber.Ctx) error {
	resp, err := resty.New().R().
		SetHeaders(Headers{
			"Authorization": "Bearer " + p.Access,
			"Accept":        "application/json",
		}).
		Get(SpotifyAPI + p.Endpoint)
	if err != nil {
		LogError("Proxy-PlaylistRequest", "Requesting "+p.Endpoint, err)
		return InternalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 200:
		c.Set("Content-Type", "application/json")
		return c.Status(fiber.StatusOK).Send(resp.Body())
	case 400:
		return BadRequest(c, "invalid playlist-id")
	case 404:
		return BadRequest(c, "playlist not found", fiber.StatusNotFound)
	default:
		LogError(
			"Proxy-PlaylistRequest",
			"Requesting "+resp.Request.URL,
			errors.New(strconv.Itoa(resp.StatusCode())+": "+string(resp.Body())))
		return InternalServerError(c, "error requesting "+c.Path())
	}
}

// TrackRequest proxies a request to the Spotify API for track.
// returns 200 if the track is found.
// returns 400 if the track-id is invalid.
// returns 404 if the track is not found.
func (p *Proxy) TrackRequest(c *fiber.Ctx) error {
	resp, err := resty.New().R().
		SetHeaders(Headers{
			"Authorization": "Bearer " + p.Access,
			"Accept":        "application/json",
		}).
		Get(SpotifyAPI + p.Endpoint)
	if err != nil {
		LogError("Proxy-TrackRequest", "Requesting "+p.Endpoint, err)
		return InternalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 200:
		c.Set("Content-Type", "application/json")
		return c.Status(fiber.StatusOK).Send(resp.Body())
	case 400:
		return BadRequest(c, "invalid track-id")
	case 404:
		return BadRequest(c, "track not found", fiber.StatusBadRequest)
	default:
		LogError(
			"Proxy-TrackRequest",
			"Requesting "+c.Path(),
			errors.New(strconv.Itoa(resp.StatusCode())+": "+string(resp.Body())),
		)
		return InternalServerError(c, "error requesting "+c.Path())
	}
}
