package handlers

import (
	"github.com/gofiber/fiber/v2"
	. "groove/pkgs/util"
	"strconv"
)

func (h *Handlers) Search(c *fiber.Ctx) error {
	queryTypes := c.Query("type")
	if queryTypes == "" {
		return BadRequest(c, "type is required")
	}

	return h.actions.Search(c,
		c.Params("query"),
		queryTypes,
		c.Query("market", "US"),
		strconv.Itoa(c.QueryInt("limit", 18)),
	)
}
