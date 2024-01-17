package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) GetTrack(c *fiber.Ctx) error {
	return h.actions.GetTrack(c, c.Params("id"))
}
