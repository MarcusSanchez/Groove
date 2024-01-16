package handlers

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (h *Handlers) Search(c *fiber.Ctx) error {

	if c.Query("type") == "" {
		return c.Status(400).SendString("invalid 'type' query parameter")
	}

	response := h.actions.Search(c,
		c.Params("query"),
		c.Query("type"),
		c.Query("market", "US"),
		strconv.Itoa(
			c.QueryInt("limit", 18),
		),
	)
	return response
}
