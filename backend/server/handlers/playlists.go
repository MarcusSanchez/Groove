package handlers

import (
	"github.com/gofiber/fiber/v2"
	. "groove/pkgs/util"
	"strconv"
)

func (h *Handlers) GetAllPlaylists(c *fiber.Ctx) error {
	return h.Actions.GetAllPlaylists(c)
}

func (h *Handlers) GetPlaylistWithTracks(c *fiber.Ctx) error {
	return h.Actions.GetPlaylist(c, c.Params("id"))
}

func (h *Handlers) GetMorePlaylistTracks(c *fiber.Ctx) error {
	offset := c.QueryInt("offset")
	if offset == 0 {
		return BadRequest(c, "invalid offset")
	}

	return h.Actions.GetMorePlaylistTracks(c, c.Params("id"), strconv.Itoa(offset))
}

func (h *Handlers) AddTrackToPlaylist(c *fiber.Ctx) error {
	trackID := c.Query("id")
	if trackID == "" {
		return BadRequest(c, "track-id is required")
	}

	return h.Actions.AddTrackToPlaylist(c, c.Params("id"), trackID)
}

func (h *Handlers) RemoveTrackFromPlaylist(c *fiber.Ctx) error {
	trackID := c.Query("id")
	if trackID == "" {
		return BadRequest(c, "track-id is required")
	}

	return h.Actions.RemoveTrackFromPlaylist(c, c.Params("id"), trackID)
}
