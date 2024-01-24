package entgo

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"groove/pkgs/env"
	. "groove/pkgs/util"
	"groove/services/database/entgo/ent"
	"groove/types"
)

var _ types.DatabaseService = (*Database)(nil)

type Database struct {
	client *ent.Client

	users        types.UserService
	sessions     types.SessionService
	oauthStates  types.OAuthStateService
	spotifylinks types.SpotifyLinkService
}

func NewDatabase(lc fx.Lifecycle, shutdowner fx.Shutdowner, env *env.Env) types.DatabaseService {
	database := &Database{}
	if err := database.Open(env.PgURI); err != nil {
		LogError("entgo-NewDatabase", "failed connecting to postgresql", err)
		_ = shutdowner.Shutdown()
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return database.Migrate(ctx)
		},
		OnStop: func(context.Context) error {
			return database.Close()
		},
	})

	database.users = NewUserService(database.client)

	return database
}

func (e *Database) Open(uri string) (err error) {
	e.client, err = ent.Open("postgres", uri)
	return err
}

func (e *Database) Close() error {
	return e.client.Close()
}

func (e *Database) Migrate(ctx context.Context) error {
	if err := e.client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed creating schema resources: %w", err)
	}
	return nil
}

func (e *Database) Users() types.UserService {
	return e.users
}
