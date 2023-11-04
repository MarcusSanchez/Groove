package middleware

import (
	"GrooveGuru/env"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recovery "github.com/gofiber/fiber/v2/middleware/recover"
)

func Attach(app *fiber.App) {
	app.Use(logger.New())
	app.Use(recovery.New())
	if !env.IsProd {
		// in development, frontend and backend are listening on different ports;
		// therefore CORS needs to be configured to allow all origins.
		app.Use(cors.New())
	}
}
