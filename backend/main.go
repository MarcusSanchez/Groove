package main

import (
	"GrooveGuru/pkgs/db"
	"GrooveGuru/pkgs/env"
	"GrooveGuru/server"
	"GrooveGuru/server/actions"
	"GrooveGuru/server/handlers"
	"GrooveGuru/server/middleware"

	_ "github.com/lib/pq"
	"go.uber.org/fx"
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
