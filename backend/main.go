package main

import (
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"groove/pkgs/db"
	"groove/pkgs/env"
	"groove/server"
)

func main() {
	fx.New(
		fx.Provide(
			db.ProvideClient,
			env.ProvideEnvVars,
		),
		fx.Invoke(
			server.InvokeServer,
			db.InvokeScheduler,
		),
	).Run()
}
