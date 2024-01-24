package actions

import (
	"groove/pkgs/ent"
	"groove/pkgs/env"
)

type Actions struct {
	Client *ent.Client
	Env    *env.Env
}
