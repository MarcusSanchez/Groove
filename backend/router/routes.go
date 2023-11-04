package router

import "github.com/gofiber/fiber/v2"

func Start(app *fiber.App) {
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("OK")
	})
}
