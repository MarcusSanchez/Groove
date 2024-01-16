package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) GetAlbum(c *fiber.Ctx) error {
	albumID := c.Params("id")
	if albumID == "" {
		return c.Status(400).SendString("invalid album-id")
	}

	response := h.actions.GetAlbum(c, albumID)
	return response
}

func (h *Handlers) GetAlbumTracks(c *fiber.Ctx) error {
	albumID := c.Params("id")
	if albumID == "" {
		return c.Status(400).SendString("invalid album-id")
	}

	response := h.actions.GetAlbumTracks(c, albumID)
	return response
}
