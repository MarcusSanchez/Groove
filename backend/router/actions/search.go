package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

func Search(c *fiber.Ctx, query, Type, market string, limit string) error {
	access := c.Locals("access").(string)

	qParams := urlSearchParams(params{
		"q":      query,
		"type":   Type,
		"market": market,
		"limit":  limit,
	})

	resp, err := resty.New().R().
		SetHeaders(headers{
			"Authorization": "Bearer " + access,
			"Accept":        "application/json",
		}).
		Get("https://api.spotify.com/v1/search?" + qParams)
	if err != nil {
		logError("Search", "Requesting ", err)
		return internalServerError(c, "error requesting "+c.Path())
	}

	switch resp.StatusCode() {
	case 400:
		return badRequest(c, "invalid search")
	case 200:
		var data map[string]any
		_ = json.Unmarshal(resp.Body(), &data)
		return c.Status(200).JSON(data)
	default:
		logError(
			"Search",
			"Requesting "+resp.Request.URL,
			errors.New(fmt.Sprintln(resp.StatusCode(), ", ", string(resp.Body()))),
		)
		return internalServerError(c, "error requesting "+c.Path())
	}
}
