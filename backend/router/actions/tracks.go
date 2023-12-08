package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

// GetTrack returns a track object from the Spotify API.
// returns 400 if the track-id is invalid.
// returns 404 if the track is not found.
// returns 200 if the track is found.
func GetTrack(c *fiber.Ctx, trackID string) error {
	access := c.Locals("access").(string)

	spotify, response := trackRequest(c,
		"/tracks/"+trackID+"?market=US",
		access,
	)
	if spotify == nil {
		return response
	}

	return trackResponse(c, spotify)
}

/** helpers **/
func trackRequest(c *fiber.Ctx, endpoint, access string) (*resty.Response, error) {
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

func trackResponse(c *fiber.Ctx, resp *resty.Response) error {
	switch resp.StatusCode() {
	case 400:
		return badRequest(c, "invalid track-id")
	case 404:
		return badRequest(c, "track not found", 404)
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
