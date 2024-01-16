package actions

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// GetAllPlaylists returns all playlists for the current user.
// returns 200 on success.
func (*Actions) GetAllPlaylists(c *fiber.Ctx) error {
	access := c.Locals("access").(string)

	spotify, response := playlistRequest(c,
		"/me/playlists?limit=50",
		access,
	)
	if spotify == nil {
		return response
	}

	return playlistResponse(c, spotify)
}

// GetPlaylist returns a playlist with the first 100 tracks with the given id.
// returns 200 on success.
// returns 404 if the playlist is not found.
// returns 400 if the playlist-id is invalid.
func (*Actions) GetPlaylist(c *fiber.Ctx, playlistID string) error {
	access := c.Locals("access").(string)

	spotify, response := playlistRequest(c,
		"/playlists/"+playlistID+"?market=US&limit=100",
		access,
	)
	if spotify == nil {
		return response
	}

	return playlistResponse(c, spotify)
}

// GetMorePlaylistTracks returns a playlist with the next 100 tracks with the given id.
// returns 200 on success.
func (*Actions) GetMorePlaylistTracks(c *fiber.Ctx, playlistID string, offset int) error {
	access := c.Locals("access").(string)

	spotify, response := playlistRequest(c,
		"/playlists/"+playlistID+"/tracks?market=US&limit=100&offset="+strconv.Itoa(offset),
		access,
	)
	if spotify == nil {
		return response
	}

	return playlistResponse(c, spotify)
}

// AddTrackToPlaylist adds a track to a playlist with the given ids.
// returns 201 on success.
// returns 404 if the playlist is not found.
// returns 403 if the playlist is not collaborative.
// returns 400 if the playlist-id is invalid.
func (*Actions) AddTrackToPlaylist(c *fiber.Ctx, playlistID, trackID string) error {
	access := c.Locals("access").(string)

	endpoint := "/playlists/" + playlistID + "/tracks"
	http := resty.New()
	resp, err := http.R().
		SetHeaders(headers{
			"Authorization": "Bearer " + access,
			"Accept":        "application/json",
		}).
		SetBody(`{"uris":["spotify:track:` + trackID + `"]}`).
		Post("https://api.spotify.com/v1" + endpoint)
	if err != nil {
		logError("AddTrackToPlaylist", "Requesting "+endpoint, err)
		return internalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 201:
		return c.Status(201).SendString("track added to playlist")
	case 400:
		return badRequest(c, "invalid track-id")
	case 403:
		return badRequest(c, "playlist is not collaborative")
	case 404:
		return badRequest(c, "playlist not found", 404)
	default:
		logError(
			"AddTrackToPlaylist",
			"Requesting "+resp.Request.URL,
			errors.New(strconv.Itoa(resp.StatusCode())+", "+string(resp.Body())),
		)
		return internalServerError(c, "error requesting "+c.Path())
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

	http := resty.New()
	resp, err := http.R().
		SetHeaders(headers{
			"Authorization": "Bearer " + access,
			"Accept":        "application/json",
		}).
		SetBody(`{"tracks":[{"uri":"spotify:track:` + trackID + `"}]}`).
		Delete("https://api.spotify.com/v1" + endpoint)
	if err != nil {
		logError("artistRequest", "Requesting "+endpoint, err)
		return internalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 200:
		return c.Status(200).SendString("track removed from playlist")
	case 400:
		return badRequest(c, "invalid track-id")
	case 403:
		return badRequest(c, "playlist is not collaborative")
	case 404:
		return badRequest(c, "playlist not found", 404)
	default:
		logError(
			"RemoveTrackFromPlaylist",
			"Requesting "+resp.Request.URL,
			errors.New(strconv.Itoa(resp.StatusCode())+", "+string(resp.Body())),
		)
		return internalServerError(c, "error requesting "+c.Path())
	}
}

/** helpers **/

// playlistRequest proxy request to the spotify api with the given endpoint and access token.
func playlistRequest(c *fiber.Ctx, endpoint, access string) (*resty.Response, error) {
	resp, err := resty.New().R().
		SetHeaders(headers{
			"Authorization": "Bearer " + access,
			"Accept":        "application/json",
		}).
		Get("https://api.spotify.com/v1" + endpoint)
	if err != nil {
		logError("artistRequest", "Requesting "+endpoint, err)
		return nil, internalServerError(c, "error requesting "+c.Path())
	}

	return resp, nil
}

// playlistResponse handles the response from the spotify api.
func playlistResponse(c *fiber.Ctx, resp *resty.Response) error {
	switch resp.StatusCode() {
	case 400:
		return badRequest(c, "invalid playlist-id")
	case 404:
		return badRequest(c, "playlist not found", 404)
	case 200:
		var data map[string]any
		_ = json.Unmarshal(resp.Body(), &data)
		return c.Status(200).JSON(data)
	default:
		logError(
			"playlistResponse",
			"Requesting "+resp.Request.URL,
			errors.New(strconv.Itoa(resp.StatusCode())+", "+string(resp.Body())),
		)
		return internalServerError(c, "error requesting "+c.Path())
	}
}
