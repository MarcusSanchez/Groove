package db

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"groove/pkgs/ent"
	"groove/pkgs/env"
	. "groove/pkgs/util"
)

func ProvideClient(lc fx.Lifecycle, shutdowner fx.Shutdowner, env *env.Env) *ent.Client {
	client, err := ent.Open("postgres", env.PgURI)
	if err != nil {
		LogError("ProvideClient", "failed connecting to postgresql", err)
		_ = shutdowner.Shutdown()
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
