package handlers

import (
	"GrooveGuru/router/actions"
	"github.com/gofiber/fiber/v2"
)

func GetArtist(c *fiber.Ctx) error {
	artistID := c.Params("id")
	if artistID == "" {
		return c.Status(400).SendString("invalid artist-id")
	}

	response := actions.GetArtist(c, artistID)
	return response
}

func GetRelatedArtists(c *fiber.Ctx) error {
	artistID := c.Params("id")
	if artistID == "" {
		return c.Status(400).SendString("invalid artist-id")
	}

	response := actions.GetRelatedArtists(c, artistID)
	return response

}

func GetArtistTopTracks(c *fiber.Ctx) error {
	artistID := c.Params("id")
	if artistID == "" {
		return c.Status(400).SendString("invalid artist-id")
	}

	response := actions.GetArtistTopTracks(c, artistID)
	return response
}

func GetArtistAlbums(c *fiber.Ctx) error {
	artistID := c.Params("id")
	if artistID == "" {
		return c.Status(400).SendString("invalid artist-id")
	}

	response := actions.GetArtistAlbums(c, artistID)
	return response
}
