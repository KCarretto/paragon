package main

import (
	"context"

	"github.com/kcarretto/paragon/ent"
	_ "github.com/mattn/go-sqlite3"
)

func getClient(ctx context.Context) *ent.Client {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		panic(err)
	}
	if err = client.Schema.Create(ctx); err != nil {
		panic(err)
	}
	return client
}
