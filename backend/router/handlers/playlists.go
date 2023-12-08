package handlers

import (
	"GrooveGuru/router/actions"
	"github.com/gofiber/fiber/v2"
)

func GetAllPlaylists(c *fiber.Ctx) error {
	response := actions.GetAllPlaylists(c)
	return response
}

func GetPlaylistWithTracks(c *fiber.Ctx) error {
	playlistID := c.Params("id")
	if playlistID == "" {
		return c.Status(400).SendString("invalid playlist-id")
	}

	response := actions.GetPlaylist(c, playlistID)
	return response
}

func GetMorePlaylistTracks(c *fiber.Ctx) error {
	playlistID := c.Params("id")
	if playlistID == "" {
		return c.Status(400).SendString("invalid playlist-id")
	}

	offset := c.QueryInt("offset")
	if offset == 0 {
		return c.Status(400).SendString("invalid offset")
	}

	response := actions.GetMorePlaylistTracks(c, playlistID, offset)
	return response
}

func AddTrackToPlaylist(c *fiber.Ctx) error {
	playlistID := c.Params("id")
	if playlistID == "" {
		return c.Status(400).SendString("invalid playlist-id")
	}

	trackID := c.Query("id")
	if trackID == "" {
		return c.Status(400).SendString("invalid track-id")
	}

	response := actions.AddTrackToPlaylist(c, playlistID, trackID)
	return response
}

func RemoveTrackFromPlaylist(c *fiber.Ctx) error {
	playlistID := c.Params("id")
	if playlistID == "" {
		return c.Status(400).SendString("invalid playlist-id")
	}

	trackID := c.Query("id")
	if trackID == "" {
		return c.Status(400).SendString("invalid track-id")
	}

	response := actions.RemoveTrackFromPlaylist(c, playlistID, trackID)
	return response
}
