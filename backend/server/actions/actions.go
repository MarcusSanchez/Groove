package actions

import (
	"groove/pkgs/ent"
	"groove/pkgs/env"
)

type Actions struct {
	client *ent.Client
	env    *env.Env
}

func ProvideActions(client *ent.Client, env *env.Env) *Actions {
	return &Actions{
		client: client,
		env:    env,
	}
}
