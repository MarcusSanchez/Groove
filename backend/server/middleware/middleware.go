package middleware

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recovery "github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
	"groove/pkgs/ent"
	SpotifyLink "groove/pkgs/ent/spotifylink"
	"groove/pkgs/env"
	. "groove/pkgs/util"
)

type Middlewares struct {
	Client     *ent.Client
	Env        *env.Env
	Shutdowner fx.Shutdowner
}

// Attach attaches the middleware that run on all endpoints.
func (m *Middlewares) Attach(app *fiber.App) {
	app.Static("/", "./public")
	app.Use(recovery.New()) /* recovers from panics and delivers an internal server error */
	app.Use(logger.New())
	app.Use(m.ReactServer) /* serves the frontend */
	switch m.Env.IsProd {
	case false:
		// in development, frontend and backend are listening on different ports;
		// therefore CORS needs to be configured to allow the frontend url.
		app.Use(cors.New(cors.Config{
			AllowOrigins:     m.Env.FrontendURL,
			AllowCredentials: true,
		}))
	}
}

// ReactServer serves the frontend.
// this is used for the catch-all route.
// if route starts with /api, it will not be served by this function.
func (*Middlewares) ReactServer(c *fiber.Ctx) error {
	path := c.Path()
	if len(path) > 4 && path[:4] == "/api" {
		return c.Next()
	}
	return c.SendFile("./public/index.html")
}

// defaultAccessToken utility function that returns the access token of the default user.
func (m *Middlewares) defaultAccessToken() (*ent.SpotifyLink, error) {
	link, err := m.Client.SpotifyLink.
		Query().
		Where(SpotifyLink.UserIDEQ(1)).
		First(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			LogError("defaultAccessToken", "default user not set or deleted", err)
			_ = m.Shutdowner.Shutdown()
		}
		return nil, err
	}

	return link, nil
}
