package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

// GetArtist returns the artist with the given id.
// Returns 400 if the artist-id is invalid.
// Returns 404 if the artist is not found.
// Returns 200 with the artist data if successful.
func (*Actions) GetArtist(c *fiber.Ctx, artistID string) error {
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

// GetRelatedArtists returns the artists related to the artist with the given id.
// Returns 400 if the artist-id is invalid.
// Returns 404 if the artist is not found.
// Returns 200 with the artist data if successful.
func (*Actions) GetRelatedArtists(c *fiber.Ctx, artistID string) error {
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

// GetArtistTopTracks returns the top tracks of the artist with the given id.
// Returns 400 if the artist-id is invalid.
// Returns 404 if the artist is not found.
// Returns 200 with the artist data if successful.
func (*Actions) GetArtistTopTracks(c *fiber.Ctx, artistID string) error {
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

func (*Actions) GetArtistAlbums(c *fiber.Ctx, artistID string) error {
	access := c.Locals("access").(string)

	spotify, response := artistRequest(c,
		"/artists/"+artistID+"/albums?market=US&limit=50&include_groups=album",
		access,
	)
	if spotify == nil {
		return response
	}

	return artistResponse(c, spotify)
}

/** helpers **/

// artistRequest is a Proxy for the Spotify API's artist endpoints.
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

// artistResponse handles the response from the Spotify API.
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
