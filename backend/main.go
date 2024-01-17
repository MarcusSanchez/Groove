package main

import (
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"groove/pkgs/db"
	"groove/pkgs/env"
	"groove/server"
	"groove/server/actions"
	"groove/server/handlers"
	"groove/server/middleware"
)

func main() {
	fx.New(
		fx.Provide(
			db.ProvideClient,
			env.ProvideEnvVars,
			middleware.ProvideMiddlewares,
			handlers.ProvideHandlers,
			actions.ProvideActions,
		),
		fx.Invoke(
			server.InvokeFiber,
			db.InvokeScheduler,
		),
	).Run()
}
