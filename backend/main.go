package main

import (
	"GrooveGuru/db"
	"GrooveGuru/env"
	"GrooveGuru/router"
	"GrooveGuru/router/middleware"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	db.Migrations()
	defer db.Close()

	app := fiber.New()
	middleware.Attach(app)
	router.Start(app)

	_ = app.Listen(":" + env.Port)
}
