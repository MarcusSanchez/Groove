package handlers

import "github.com/gofiber/fiber/v2"

// ReactServer serves the frontend.
func ReactServer(c *fiber.Ctx) error {
	return c.SendFile("./public/index.html")
}

// Health returns a 200 if the server is running.
func Health(c *fiber.Ctx) error {
	return c.Status(200).SendString("OK")
}
