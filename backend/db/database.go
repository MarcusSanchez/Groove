package db

import (
	"GrooveGuru/ent"
	"GrooveGuru/env"
	"context"
	_ "github.com/lib/pq"
	"log"
)

var ctx = context.Background()
var client *ent.Client

func init() {
	var err error
	client, err = ent.Open("postgres", env.PgURI)
	if err != nil {
		log.Fatal("failed connecting to postgresql: ", err)
	}
}

func Migrations() {
	// init will handle the majority of the work
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatal("failed creating schema resources: ", err)
	}
}

func Instance() (*ent.Client, context.Context) {
	return client, ctx
}

func Close() {
	if err := client.Close(); err != nil {
		log.Fatal("failed to close client: ", err)
	}
}
