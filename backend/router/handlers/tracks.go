package handlers

import (
	"GrooveGuru/router/actions"
	"github.com/gofiber/fiber/v2"
)

func GetTrack(c *fiber.Ctx) error {
	trackID := c.Params("id")
	if trackID == "" {
		return c.Status(400).SendString("invalid track-id")
	}

	response := actions.GetTrack(c, trackID)
	return response
}
