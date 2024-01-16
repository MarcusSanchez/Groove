package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) GetArtist(c *fiber.Ctx) error {
	artistID := c.Params("id")
	if artistID == "" {
		return c.Status(400).SendString("invalid artist-id")
	}

	response := h.actions.GetArtist(c, artistID)
	return response
}

func (h *Handlers) GetRelatedArtists(c *fiber.Ctx) error {
	artistID := c.Params("id")
	if artistID == "" {
		return c.Status(400).SendString("invalid artist-id")
	}

	response := h.actions.GetRelatedArtists(c, artistID)
	return response

}

func (h *Handlers) GetArtistTopTracks(c *fiber.Ctx) error {
	artistID := c.Params("id")
	if artistID == "" {
		return c.Status(400).SendString("invalid artist-id")
	}

	response := h.actions.GetArtistTopTracks(c, artistID)
	return response
}

func (h *Handlers) GetArtistAlbums(c *fiber.Ctx) error {
	artistID := c.Params("id")
	if artistID == "" {
		return c.Status(400).SendString("invalid artist-id")
	}

	response := h.actions.GetArtistAlbums(c, artistID)
	return response
}
