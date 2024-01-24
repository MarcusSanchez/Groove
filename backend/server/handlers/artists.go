package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) GetArtist(c *fiber.Ctx) error {
	return h.Actions.GetArtist(c, c.Params("id"))
}

func (h *Handlers) GetRelatedArtists(c *fiber.Ctx) error {
	return h.Actions.GetRelatedArtists(c, c.Params("id"))
}

func (h *Handlers) GetArtistTopTracks(c *fiber.Ctx) error {
	return h.Actions.GetArtistTopTracks(c, c.Params("id"))
}

func (h *Handlers) GetArtistAlbums(c *fiber.Ctx) error {
	return h.Actions.GetArtistAlbums(c, c.Params("id"))
}
