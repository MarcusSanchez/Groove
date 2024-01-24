package actions

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	. "groove/pkgs/util"
	"net/http"
	"strconv"
)

// GetAllPlaylists returns all playlists for the current user.
func (*Actions) GetAllPlaylists(c *fiber.Ctx) error {
	proxy := &Proxy{
		Endpoint: "/me/playlists?limit=50",
		Access:   c.Locals("access").(string),
	}
	return proxy.PlaylistRequest(c)
}

// GetPlaylist returns a playlist with the first 100 tracks with the given id.
func (*Actions) GetPlaylist(c *fiber.Ctx, playlistID string) error {
	proxy := &Proxy{
		Endpoint: "/playlists/" + playlistID + "?market=US&limit=100",
		Access:   c.Locals("access").(string),
	}
	return proxy.PlaylistRequest(c)
}

// GetMorePlaylistTracks returns a playlist with the next 100 tracks with the given id.
func (*Actions) GetMorePlaylistTracks(c *fiber.Ctx, playlistID, offset string) error {
	proxy := &Proxy{
		Endpoint: "/playlists/" + playlistID + "/tracks?market=US&limit=100&offset=" + offset,
		Access:   c.Locals("access").(string),
	}
	return proxy.PlaylistRequest(c)
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
		Post(SpotifyAPI + endpoint)
	if err != nil {
		LogError("AddTrackToPlaylist", "Requesting "+endpoint, err)
		return InternalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 201:
		return c.Status(http.StatusCreated).SendString("track added to playlist")
	case 400:
		return BadRequest(c, "invalid track-id")
	case 403:
		return BadRequest(c, "playlist is not collaborative")
	case 404:
		return BadRequest(c, "playlist not found", http.StatusNotFound)
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
		Delete(SpotifyAPI + endpoint)
	if err != nil {
		LogError("RemoveTrackFromPlaylist", "Requesting "+endpoint, err)
		return InternalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 200:
		return c.Status(http.StatusOK).SendString("track removed from playlist")
	case 400:
		return BadRequest(c, "invalid track-id")
	case 403:
		return BadRequest(c, "playlist is not collaborative")
	case 404:
		return BadRequest(c, "playlist not found", http.StatusNotFound)
	default:
		LogError(
			"RemoveTrackFromPlaylist",
			"Requesting "+c.Path(),
			errors.New(strconv.Itoa(resp.StatusCode())+": "+string(resp.Body())),
		)
		return InternalServerError(c, "error requesting "+c.Path())
	}
}
