// +build !dev

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	entsql "github.com/facebookincubator/ent/dialect/sql"
	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/migrate"

	_ "github.com/go-sql-driver/mysql"
	"gocloud.dev/mysql"
	// _ "gocloud.dev/mysql/gcpmysql"
)

func newGraph(ctx context.Context) *ent.Client {
	var mysqlDSN string
	if mysqlDSN = os.Getenv("PG_MYSQL_DSN"); mysqlDSN == "" {
		panic(fmt.Errorf("failed to connect to mysql: missing PG_MYSQL_DSN"))
	}

	db, err := mysql.Open(ctx, mysqlDSN)
	if err != nil {
		panic(fmt.Errorf("failed to connect to mysql: %w", err))
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB("mysql", db)

	client := ent.NewClient(ent.Driver(drv))
	if err := client.Schema.Create(ctx, migrate.WithGlobalUniqueID(true)); err != nil {
		panic(err)
	}

	return client
}
