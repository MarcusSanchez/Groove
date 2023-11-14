package handlers

import (
	"GrooveGuru/router/actions"
	"github.com/gofiber/fiber/v2"
)

func GetAlbum(c *fiber.Ctx) error {
	albumID := c.Params("id")
	if albumID == "" {
		return c.Status(400).SendString("invalid album-id")
	}

	response := actions.GetAlbum(c, albumID)
	return response
}

func GetAlbumTracks(c *fiber.Ctx) error {
	albumID := c.Params("id")
	if albumID == "" {
		return c.Status(400).SendString("invalid album-id")
	}

	response := actions.GetAlbumTracks(c, albumID)
	return response
}
