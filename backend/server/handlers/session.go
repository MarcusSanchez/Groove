package handlers

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (h *Handlers) Register(c *fiber.Ctx) error {

	type RegisterUserPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	var payload RegisterUserPayload
	if c.BodyParser(&payload) != nil {
		return c.Status(400).SendString("Invalid JSON")
	}

	response := h.actions.Register(c,
		payload.Password,
		payload.Username,
		strings.ToLower(payload.Email),
	)
	return response
}

func (h *Handlers) Login(c *fiber.Ctx) error {

	type LoginPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var payload LoginPayload
	if c.BodyParser(&payload) != nil {
		return c.Status(400).SendString("Invalid JSON")
	}

	response := h.actions.Login(c,
		payload.Username,
		payload.Password,
	)
	return response
}

func (h *Handlers) Logout(c *fiber.Ctx) error {
	response := h.actions.Logout(c)
	return response
}

func (h *Handlers) Authenticate(c *fiber.Ctx) error {
	response := h.actions.Authenticate(c)
	return response
}
