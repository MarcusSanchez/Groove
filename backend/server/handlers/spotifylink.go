package handlers

import (
	"github.com/gofiber/fiber/v2"
	. "groove/pkgs/util"
)

func (h *Handlers) LinkSpotify(c *fiber.Ctx) error {
	return h.Actions.LinkSpotify(c)
}

func (h *Handlers) SpotifyCallback(c *fiber.Ctx) error {
	code, state := c.Query("code"), c.Query("state")
	if code == "" || state == "" {
		return BadRequest(c, "invalid query params")
	}

	return h.Actions.SpotifyCallback(c, code, state)
}

func (h *Handlers) UnlinkSpotify(c *fiber.Ctx) error {
	return h.Actions.UnlinkSpotify(c)
}

func (h *Handlers) GetCurrentUser(c *fiber.Ctx) error {
	return h.Actions.GetCurrentUser(c)
}
