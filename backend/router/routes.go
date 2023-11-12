package router

import (
	"GrooveGuru/router/handlers"
	"GrooveGuru/router/middleware"
	"github.com/gofiber/fiber/v2"
)

func Start(app fiber.Router) {
	/** api endpoints **/
	api := app.Group("/api")
	api.Get("/health", handlers.Health)

	/** session endpoints **/
	api.Post("/register", middleware.RedirectAuthorized, handlers.Register)
	api.Post("/login", middleware.RedirectAuthorized, handlers.Login)
	api.Post("/logout", middleware.CheckCSRF, handlers.Logout)
	api.Post("/authenticate", middleware.CheckCSRF, handlers.Authenticate)

	/** spotify-link endpoints **/
	api.Post("/spotify/link", middleware.CheckCSRF, middleware.RedirectLinked, handlers.SpotifyLink)
	api.Get("/spotify/callback", middleware.AuthorizeAny, handlers.SpotifyCallback)

	// catch-all route for the frontend.
	// placed after all other routes to prevent conflicts.
	app.Get("*", handlers.ReactServer)
}
