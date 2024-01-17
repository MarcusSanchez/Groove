package server

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"groove/pkgs/env"
	. "groove/pkgs/util"
	"groove/server/handlers"
	"groove/server/middleware"
)

func InvokeFiber(lc fx.Lifecycle, shutdowner fx.Shutdowner, handlers *handlers.Handlers, mw *middleware.Middlewares, env *env.Env) {
	app := fiber.New()
	mw.Attach(app)
	SetupEndpoints(app, handlers, mw)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := app.Listen(env.Port); err != nil {
					LogError("InvokeFiber", "failed to listen", err)
					_ = shutdowner.Shutdown()
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			return app.Shutdown()
		},
	})
}
