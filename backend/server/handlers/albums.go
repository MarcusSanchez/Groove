package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) GetAlbum(c *fiber.Ctx) error {
	return h.actions.GetAlbum(c, c.Params("id"))
}

func (h *Handlers) GetAlbumTracks(c *fiber.Ctx) error {
	return h.actions.GetAlbumTracks(c, c.Params("id"))
}
