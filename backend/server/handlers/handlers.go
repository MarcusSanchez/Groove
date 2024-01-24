package handlers

import (
	"github.com/gofiber/fiber/v2"
	"groove/server/actions"
	"net/http"
)

type Handlers struct {
	Actions *actions.Actions
}

// Health returns a 200 if the server is running.
func (*Handlers) Health(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).SendString("OK")
}
