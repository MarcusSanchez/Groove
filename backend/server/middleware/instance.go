package middleware

import (
	"GrooveGuru/pkgs/ent"
	"GrooveGuru/pkgs/env"
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
