package router

import (
	"GrooveGuru/middleware"
	"GrooveGuru/router/handlers"
	"github.com/gofiber/fiber/v2"
)

func Start(app fiber.Router) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("OK")
	})

	app.Post("/register", handlers.Register).Use(middleware.RedirectAuthorized)
	app.Post("/login", handlers.Login).Use(middleware.RedirectAuthorized)
	app.Post("/logout", handlers.Logout).Use(middleware.CheckCSRF)
}
