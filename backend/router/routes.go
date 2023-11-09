package router

import (
	"GrooveGuru/router/handlers"
	"GrooveGuru/router/middleware"
	"github.com/gofiber/fiber/v2"
)

func Start(app fiber.Router) {
	// catch-all route for the frontend.
	app.Get("*", handlers.ReactServer)

	/** api endpoints **/
	api := app.Group("/api")
	api.Get("/health", handlers.Health)

	/** session endpoints **/
	api.Post("/register", middleware.RedirectAuthorized, handlers.Register)
	api.Post("/login", middleware.RedirectAuthorized, handlers.Login)
	api.Post("/logout", middleware.CheckCSRF, handlers.Logout)
	api.Post("/Authenticate", middleware.CheckCSRF, handlers.Authenticate)

}
