package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"groove/pkgs/env"
	"groove/services/database/entgo"
	"groove/types"
)

func main() {
	fx.New(
		fx.Provide(
			env.ProvideEnvVars,
			entgo.NewDatabase,
		),
		fx.Invoke(
			FirstUser,
		),
	).Run()
}

func FirstUser(db types.DatabaseService) {
	u, err := db.Users().FindByID(1, context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("Username:", u.Username)
	fmt.Println("Password:", u.Password)
	fmt.Println("Email:", u.Email)
}
