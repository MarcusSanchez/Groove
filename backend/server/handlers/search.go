package handlers

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (h *Handlers) Search(c *fiber.Ctx) error {
	queryTypes := c.Query("type")
	if queryTypes == "" {
		return c.Status(400).SendString("type is required")
	}

	return h.actions.Search(c,
		c.Params("query"),
		queryTypes,
		c.Query("market", "US"),
		strconv.Itoa(c.QueryInt("limit", 18)),
	)
}
