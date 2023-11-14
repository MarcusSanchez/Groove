package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

func GetAlbum(c *fiber.Ctx, albumID string) error {
	access := c.Locals("access").(string)

	spotify, response := albumRequest(c,
		"/albums/"+albumID+"?market=US",
		access,
	)
	if spotify == nil {
		return response
	}

	return albumResponse(c, spotify)
}

func GetAlbumTracks(c *fiber.Ctx, albumID string) error {
	access := c.Locals("access").(string)

	spotify, response := albumRequest(c,
		"/albums/"+albumID+"/tracks?limit=50&market=US",
		access,
	)
	if spotify == nil {
		return response
	}

	return albumResponse(c, spotify)
}

/** helpers **/

func albumRequest(c *fiber.Ctx, endpoint, access string) (*resty.Response, error) {
	resp, err := resty.New().R().
		SetHeaders(headers{
			"Authorization": "Bearer " + access,
			"Accept":        "application/json",
		}).
		Get("https://api.spotify.com/v1" + endpoint)
	if err != nil {
		logError("trackRequest", "Requesting "+endpoint, err)
		return nil, internalServerError(c, "error requesting "+c.Path())
	}

	return resp, nil
}

func albumResponse(c *fiber.Ctx, resp *resty.Response) error {
	switch resp.StatusCode() {
	case 400:
		return badRequest(c, "invalid album-id")
	case 404:
		return badRequest(c, "album not found", 404)
	case 200:
		var data map[string]any
		_ = json.Unmarshal(resp.Body(), &data)
		return c.Status(200).JSON(data)
	default:
		logError(
			"trackResponse",
			"Requesting "+resp.Request.URL,
			errors.New(fmt.Sprintln(resp.StatusCode(), ", ", string(resp.Body()))),
		)
		return internalServerError(c, "error requesting "+c.Path())
	}
}
