package server

import (
	"GrooveGuru/server/handlers"
	"GrooveGuru/server/middleware"
	"context"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"log"
)

func InvokeFiber(lc fx.Lifecycle, shutdowner fx.Shutdowner, handlers *handlers.Handlers, mw *middleware.Middlewares) {
	app := fiber.New()
	mw.Attach(app)
	SetupEndpoints(app, handlers, mw)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := app.Listen(":3000"); err != nil {
					log.Println("failed OnStart for InvokeFiber: ", err)
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
