package handlers

import (
	"GrooveGuru/router/actions"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func Search(c *fiber.Ctx) error {

	if c.Query("type") == "" {
		return c.Status(400).SendString("invalid 'type' query parameter")
	}

	response := actions.Search(c,
		c.Params("query"),
		c.Query("type"),
		c.Query("market", "US"),
		strconv.Itoa(
			c.QueryInt("limit", 18),
		),
	)
	return response
}
