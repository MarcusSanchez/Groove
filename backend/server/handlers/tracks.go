package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) GetTrack(c *fiber.Ctx) error {
	trackID := c.Params("id")
	if trackID == "" {
		return c.Status(400).SendString("invalid track-id")
	}

	response := h.actions.GetTrack(c, trackID)
	return response
}
