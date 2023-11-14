package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

func GetArtist(c *fiber.Ctx, artistID string) error {
	access := c.Locals("access").(string)

	spotify, response := artistRequest(c,
		"/artists/"+artistID,
		access,
	)
	if spotify == nil {
		return response
	}

	return artistResponse(c, spotify)
}

func GetRelatedArtists(c *fiber.Ctx, artistID string) error {
	access := c.Locals("access").(string)

	spotify, response := artistRequest(c,
		"/artists/"+artistID+"/related-artists",
		access,
	)
	if spotify == nil {
		return response
	}

	return artistResponse(c, spotify)
}

func GetArtistTopTracks(c *fiber.Ctx, artistID string) error {
	access := c.Locals("access").(string)

	spotify, response := artistRequest(c,
		"/artists/"+artistID+"/top-tracks?market=US",
		access,
	)
	if spotify == nil {
		return response
	}

	return artistResponse(c, spotify)
}

func GetArtistAlbums(c *fiber.Ctx, artistID string) error {
	access := c.Locals("access").(string)

	spotify, response := artistRequest(c,
		"/artists/"+artistID+"/albums?market=US&limit=50",
		access,
	)
	if spotify == nil {
		return response
	}

	return artistResponse(c, spotify)
}

/** helpers **/

func artistRequest(c *fiber.Ctx, endpoint, access string) (*resty.Response, error) {
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

func artistResponse(c *fiber.Ctx, resp *resty.Response) error {
	switch resp.StatusCode() {
	case 400:
		return badRequest(c, "invalid artist-id")
	case 404:
		return badRequest(c, "artist not found", 404)
	case 200:
		var data map[string]any
		_ = json.Unmarshal(resp.Body(), &data)
		return c.Status(200).JSON(data)
	default:
		logError(
			"artistResponse",
			"Requesting "+resp.Request.URL,
			errors.New(fmt.Sprintln(resp.StatusCode(), ", ", string(resp.Body()))),
		)
		return internalServerError(c, "error requesting "+c.Path())
	}
}
