package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) LinkSpotify(c *fiber.Ctx) error {
	response := h.actions.LinkSpotify(c)
	return response
}

func (h *Handlers) SpotifyCallback(c *fiber.Ctx) error {

	type CallbackQuery struct {
		Code  string `query:"code"`
		State string `query:"state"`
	}

	var query CallbackQuery
	_ = c.QueryParser(&query)
	if query.Code == "" || query.State == "" {
		return c.Status(400).SendString("Invalid query params")
	}

	response := h.actions.SpotifyCallback(c,
		query.Code,
		query.State,
	)
	return response
}

func (h *Handlers) UnlinkSpotify(c *fiber.Ctx) error {
	response := h.actions.UnlinkSpotify(c)
	return response
}

func (h *Handlers) GetCurrentUser(c *fiber.Ctx) error {
	response := h.actions.GetCurrentUser(c)
	return response
}
