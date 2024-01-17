package db

import (
	"GrooveGuru/pkgs/ent"
	"GrooveGuru/pkgs/env"
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"log"
)

func ProvideClient(lc fx.Lifecycle, env *env.Env) *ent.Client {
	client, err := ent.Open("postgres", env.PgURI)
	if err != nil {
		log.Fatal("failed connecting to postgresql: ", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err = client.Schema.Create(ctx); err != nil {
				return fmt.Errorf("failed creating schema resources: %w", err)
			}
			return nil
		},
		OnStop: func(context.Context) error {
			return client.Close()
		},
	})

	return client
}
