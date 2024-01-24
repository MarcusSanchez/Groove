package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) GetTrack(c *fiber.Ctx) error {
	return h.Actions.GetTrack(c, c.Params("id"))
}
