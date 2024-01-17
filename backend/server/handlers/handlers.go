package handlers

import (
	"github.com/gofiber/fiber/v2"
	"groove/server/actions"
)

type Handlers struct {
	actions *actions.Actions
}

func ProvideHandlers(actions *actions.Actions) *Handlers {
	return &Handlers{
		actions: actions,
	}
}

// Health returns a 200 if the server is running.
func (*Handlers) Health(c *fiber.Ctx) error {
	return c.Status(200).SendString("OK")
}
