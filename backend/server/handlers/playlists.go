package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) GetAllPlaylists(c *fiber.Ctx) error {
	response := h.actions.GetAllPlaylists(c)
	return response
}

func (h *Handlers) GetPlaylistWithTracks(c *fiber.Ctx) error {
	playlistID := c.Params("id")
	if playlistID == "" {
		return c.Status(400).SendString("invalid playlist-id")
	}

	response := h.actions.GetPlaylist(c, playlistID)
	return response
}

func (h *Handlers) GetMorePlaylistTracks(c *fiber.Ctx) error {
	playlistID := c.Params("id")
	if playlistID == "" {
		return c.Status(400).SendString("invalid playlist-id")
	}

	offset := c.QueryInt("offset")
	if offset == 0 {
		return c.Status(400).SendString("invalid offset")
	}

	response := h.actions.GetMorePlaylistTracks(c, playlistID, offset)
	return response
}

func (h *Handlers) AddTrackToPlaylist(c *fiber.Ctx) error {
	playlistID := c.Params("id")
	if playlistID == "" {
		return c.Status(400).SendString("invalid playlist-id")
	}

	trackID := c.Query("id")
	if trackID == "" {
		return c.Status(400).SendString("invalid track-id")
	}

	response := h.actions.AddTrackToPlaylist(c, playlistID, trackID)
	return response
}

func (h *Handlers) RemoveTrackFromPlaylist(c *fiber.Ctx) error {
	playlistID := c.Params("id")
	if playlistID == "" {
		return c.Status(400).SendString("invalid playlist-id")
	}

	trackID := c.Query("id")
	if trackID == "" {
		return c.Status(400).SendString("invalid track-id")
	}

	response := h.actions.RemoveTrackFromPlaylist(c, playlistID, trackID)
	return response
}
