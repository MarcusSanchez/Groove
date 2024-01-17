package actions

import (
	. "GrooveGuru/pkgs/util"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// GetAllPlaylists returns all playlists for the current user.
func (*Actions) GetAllPlaylists(c *fiber.Ctx) error {
	return proxyPlaylistRequest(c, Proxy{
		Endpoint: "/me/playlists?limit=50",
		Access:   c.Locals("access").(string),
	})
}

// GetPlaylist returns a playlist with the first 100 tracks with the given id.
func (*Actions) GetPlaylist(c *fiber.Ctx, playlistID string) error {
	return proxyPlaylistRequest(c, Proxy{
		Endpoint: "/playlists/" + playlistID + "?market=US&limit=100",
		Access:   c.Locals("access").(string),
	})
}

// GetMorePlaylistTracks returns a playlist with the next 100 tracks with the given id.
func (*Actions) GetMorePlaylistTracks(c *fiber.Ctx, playlistID, offset string) error {
	return proxyPlaylistRequest(c, Proxy{
		Endpoint: "/playlists/" + playlistID + "?market=US&limit=100&offset=" + offset,
		Access:   c.Locals("access").(string),
	})
}

// AddTrackToPlaylist adds a track to a playlist with the given ids.
// returns 201 on success.
// returns 404 if the playlist is not found.
// returns 403 if the playlist is not collaborative.
// returns 400 if the playlist-id is invalid.
func (*Actions) AddTrackToPlaylist(c *fiber.Ctx, playlistID, trackID string) error {
	access := c.Locals("access").(string)
	endpoint := "/playlists/" + playlistID + "/tracks"

	resp, err := resty.New().R().
		SetHeaders(Headers{
			"Authorization": "Bearer " + access,
			"Accept":        "application/json",
		}).
		SetBody(`{"uris":["spotify:track:` + trackID + `"]}`).
		Post("https://api.spotify.com/v1" + endpoint)
	if err != nil {
		LogError("AddTrackToPlaylist", "Requesting "+endpoint, err)
		return InternalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 201:
		return c.Status(201).SendString("track added to playlist")
	case 400:
		return BadRequest(c, "invalid track-id")
	case 403:
		return BadRequest(c, "playlist is not collaborative")
	case 404:
		return BadRequest(c, "playlist not found", 404)
	default:
		LogError(
			"AddTrackToPlaylist",
			"Requesting "+c.Path(),
			errors.New(strconv.Itoa(resp.StatusCode())+": "+string(resp.Body())),
		)
		return InternalServerError(c, "error requesting "+c.Path())
	}
}

// RemoveTrackFromPlaylist removes a track from a playlist with the given ids.
// returns 200 on success.
// returns 404 if the playlist is not found.
// returns 403 if the playlist is not collaborative.
// returns 400 if the playlist-id is invalid.
func (*Actions) RemoveTrackFromPlaylist(c *fiber.Ctx, playlistID, trackID string) error {
	access := c.Locals("access").(string)
	endpoint := "/playlists/" + playlistID + "/tracks"

	resp, err := resty.New().R().
		SetHeaders(Headers{
			"Authorization": "Bearer " + access,
			"Accept":        "application/json",
		}).
		SetBody(`{"tracks":[{ "uri":"spotify:track:` + trackID + `"}]}`).
		Delete("https://api.spotify.com/v1" + endpoint)
	if err != nil {
		LogError("RemoveTrackFromPlaylist", "Requesting "+endpoint, err)
		return InternalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 200:
		return c.Status(fiber.StatusOK).SendString("track removed from playlist")
	case 400:
		return BadRequest(c, "invalid track-id")
	case 403:
		return BadRequest(c, "playlist is not collaborative")
	case 404:
		return BadRequest(c, "playlist not found", 404)
	default:
		LogError(
			"RemoveTrackFromPlaylist",
			"Requesting "+c.Path(),
			errors.New(strconv.Itoa(resp.StatusCode())+": "+string(resp.Body())),
		)
		return InternalServerError(c, "error requesting "+c.Path())
	}
}

/* helpers **/

// proxyPlaylistRequest proxies a request to the Spotify API for a playlist.
// returns 200 with the playlist tracks if successful.
// returns 400 if the playlist-id is invalid.
// returns 404 if the playlist is not found.
func proxyPlaylistRequest(c *fiber.Ctx, proxy Proxy) error {
	resp, err := resty.New().R().
		SetHeaders(Headers{
			"Authorization": "Bearer " + proxy.Access,
			"Accept":        "application/json",
		}).
		Get("https://api.spotify.com/v1" + proxy.Endpoint)
	if err != nil {
		LogError("proxyPlaylistRequest", "Requesting "+proxy.Endpoint, err)
		return InternalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 200:
		c.Set("Content-Type", "application/json")
		return c.Status(fiber.StatusOK).Send(resp.Body())
	case 400:
		return BadRequest(c, "invalid playlist-id")
	case 404:
		return BadRequest(c, "playlist not found", 404)
	default:
		LogError(
			"proxyPlaylistResponse",
			"Requesting "+resp.Request.URL,
			errors.New(strconv.Itoa(resp.StatusCode())+", "+string(resp.Body())),
		)
		return InternalServerError(c, "error requesting "+c.Path())
	}
}
