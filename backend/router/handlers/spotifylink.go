package handlers

import (
	"GrooveGuru/router/actions"
	"github.com/gofiber/fiber/v2"
)

func LinkSpotify(c *fiber.Ctx) error {
	response := actions.LinkSpotify(c)
	return response
}

func SpotifyCallback(c *fiber.Ctx) error {

	type CallbackQuery struct {
		Code  string `query:"code"`
		State string `query:"state"`
	}

	var query CallbackQuery
	_ = c.QueryParser(&query)
	if query.Code == "" || query.State == "" {
		return c.Status(400).SendString("Invalid query params")
	}

	response := actions.SpotifyCallback(c,
		query.Code,
		query.State,
	)
	return response
}

func UnlinkSpotify(c *fiber.Ctx) error {
	response := actions.UnlinkSpotify(c)
	return response
}
