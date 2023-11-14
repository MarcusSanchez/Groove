package handlers

import "github.com/gofiber/fiber/v2"

// Health returns a 200 if the server is running.
func Health(c *fiber.Ctx) error {
	return c.Status(200).SendString("OK")
}
