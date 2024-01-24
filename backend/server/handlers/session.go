package handlers

import (
	"github.com/MarcusSanchez/go-parse"
	"github.com/gofiber/fiber/v2"
	. "groove/pkgs/util"
	"strings"
)

func (h *Handlers) Register(c *fiber.Ctx) error {

	type Payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	payload, err := parse.JSON[Payload](c.Body())
	if err != nil {
		return BadRequest(c, err.Error())
	}

	return h.Actions.Register(c,
		payload.Password,
		payload.Username,
		strings.ToLower(payload.Email),
	)
}

func (h *Handlers) Login(c *fiber.Ctx) error {

	type Payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	payload, err := parse.JSON[Payload](c.Body())
	if err != nil {
		return BadRequest(c, err.Error())
	}

	return h.Actions.Login(c,
		payload.Username,
		payload.Password,
	)
}

func (h *Handlers) Logout(c *fiber.Ctx) error {
	return h.Actions.Logout(c)
}

func (h *Handlers) Authenticate(c *fiber.Ctx) error {
	return h.Actions.Authenticate(c)
}
