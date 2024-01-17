package handlers

import (
	. "GrooveGuru/pkgs/util"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) LinkSpotify(c *fiber.Ctx) error {
	return h.actions.LinkSpotify(c)
}

func (h *Handlers) SpotifyCallback(c *fiber.Ctx) error {
	code, state := c.Query("code"), c.Query("state")
	if code == "" || state == "" {
		return BadRequest(c, "invalid query params")
	}

	return h.actions.SpotifyCallback(c, code, state)
}

func (h *Handlers) UnlinkSpotify(c *fiber.Ctx) error {
	return h.actions.UnlinkSpotify(c)
}

func (h *Handlers) GetCurrentUser(c *fiber.Ctx) error {
	return h.actions.GetCurrentUser(c)
}
