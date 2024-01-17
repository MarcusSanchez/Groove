package middleware

import (
	"GrooveGuru/pkgs/ent"
	SpotifyLink "GrooveGuru/pkgs/ent/spotifylink"
	"GrooveGuru/pkgs/env"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recovery "github.com/gofiber/fiber/v2/middleware/recover"
	"log"
)

type Middlewares struct {
	client *ent.Client
	env    *env.Env
}

func ProvideMiddlewares(client *ent.Client, env *env.Env) *Middlewares {
	return &Middlewares{
		client: client,
		env:    env,
	}
}

// Attach attaches the middleware that run on all endpoints.
func (m *Middlewares) Attach(app *fiber.App) {
	app.Static("/", "./public")
	// catch-all route for the frontend.
	app.Use("/", m.ReactServer)
	app.Use(logger.New())
	// if the server were to crash, this would restart the server.
	app.Use(recovery.New())
	switch m.env.IsProd {
	case false:
		// in development, frontend and backend are listening on different ports;
		// therefore CORS needs to be configured to allow the frontend url.
		app.Use(cors.New(cors.Config{
			AllowOrigins:     m.env.FrontendURL,
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

/* utility */
func (m *Middlewares) defaultAccessToken() (*ent.SpotifyLink, error) {
	link, err := m.client.SpotifyLink.
		Query().
		Where(SpotifyLink.UserIDEQ(1)).
		First(context.Background())
	if ent.IsNotFound(err) {
		log.Fatal("default access token not set")
	} else if err != nil {
		return nil, err
	}

	return link, err
}
