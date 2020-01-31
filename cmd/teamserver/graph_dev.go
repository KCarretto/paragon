// +build dev

package main

import (
	"context"

	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/migrate"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func newGraph(ctx context.Context, logger *zap.Logger) *ent.Client {
	logger.Debug("Connecting to sqlite")
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		panic(err)
	}

	if err := client.Schema.Create(ctx, migrate.WithGlobalUniqueID(true)); err != nil {
		panic(err)
	}

	return client
}
