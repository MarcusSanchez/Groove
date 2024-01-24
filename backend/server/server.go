package server

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"groove/pkgs/ent"
	"groove/pkgs/env"
	. "groove/pkgs/util"
	"groove/server/actions"
	"groove/server/handlers"
	"groove/server/middleware"
)

type Server struct {
	app        *fiber.App
	handlers   *handlers.Handlers
	middleware *middleware.Middlewares
}

func InvokeServer(lc fx.Lifecycle, shutdowner fx.Shutdowner, client *ent.Client, env *env.Env) {
	server := &Server{
		app: fiber.New(),
		handlers: &handlers.Handlers{
			Actions: &actions.Actions{
				Client: client,
				Env:    env,
			},
		},
		middleware: &middleware.Middlewares{
			Client:     client,
			Env:        env,
			Shutdowner: shutdowner,
		},
	}
	server.middleware.Attach(server.app)
	server.SetupEndpoints()

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := server.app.Listen(":" + env.Port); err != nil {
					LogError("InvokeServer", "failed to listen", err)
					_ = shutdowner.Shutdown()
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			return server.app.Shutdown()
		},
	})
}
