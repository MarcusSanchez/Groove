package actions

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	. "groove/pkgs/util"
	"net/http"
	"strconv"
)

// Search searches Spotify for a given query.
// returns a JSON response from Spotify.
// returns 200 if successful.
func (*Actions) Search(c *fiber.Ctx, query, Type, market string, limit string) error {
	access := c.Locals("access").(string)

	qParams := URLSearchParams(Params{
		"q":      query,
		"type":   Type,
		"market": market,
		"limit":  limit,
	})

	resp, err := resty.New().R().
		SetHeaders(Headers{
			"Authorization": "Bearer " + access,
			"Accept":        "application/json",
		}).
		Get(SpotifyAPI + "/search?" + qParams)
	if err != nil {
		LogError("Search", "Requesting ", err)
		return InternalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 200:
		c.Set("Content-Type", "application/json")
		return c.Status(http.StatusOK).Send(resp.Body())
	case 400:
		return BadRequest(c, "invalid search")
	default:
		LogError(
			"Search",
			"Requesting "+c.Path(),
			errors.New(strconv.Itoa(resp.StatusCode())+": "+string(resp.Body())),
		)
		return InternalServerError(c, "error requesting "+c.Path())
	}
}
