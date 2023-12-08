package actions

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetAllPlaylists(c *fiber.Ctx) error {
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

func GetPlaylist(c *fiber.Ctx, playlistID string) error {
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

func GetMorePlaylistTracks(c *fiber.Ctx, playlistID string, offset int) error {
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

func AddTrackToPlaylist(c *fiber.Ctx, playlistID, trackID string) error {
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

func RemoveTrackFromPlaylist(c *fiber.Ctx, playlistID, trackID string) error {
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
