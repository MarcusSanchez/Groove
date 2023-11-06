package router

import (
	"GrooveGuru/router/handlers"
	"GrooveGuru/router/middleware"
	"github.com/gofiber/fiber/v2"
)

func Start(app fiber.Router) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("OK")
	})

	/** session endpoints **/
	app.Post("/register", middleware.RedirectAuthorized, handlers.Register)
	app.Post("/login", middleware.RedirectAuthorized, handlers.Login)
	app.Post("/logout", middleware.CheckCSRF, handlers.Logout)

}
