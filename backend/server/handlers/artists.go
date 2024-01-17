package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) GetArtist(c *fiber.Ctx) error {
	return h.actions.GetArtist(c, c.Params("id"))
}

func (h *Handlers) GetRelatedArtists(c *fiber.Ctx) error {
	return h.actions.GetRelatedArtists(c, c.Params("id"))
}

func (h *Handlers) GetArtistTopTracks(c *fiber.Ctx) error {
	return h.actions.GetArtistTopTracks(c, c.Params("id"))
}

func (h *Handlers) GetArtistAlbums(c *fiber.Ctx) error {
	return h.actions.GetArtistAlbums(c, c.Params("id"))
}
