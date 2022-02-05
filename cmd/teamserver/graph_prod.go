//go:build !dev
// +build !dev

package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/kcarretto/paragon/ent"
	"github.com/kcarretto/paragon/ent/migrate"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gocloud.dev/mysql"
	_ "gocloud.dev/mysql/gcpmysql"
)

func connect(ctx context.Context, logger *zap.Logger, dsn string) (db *sql.DB, err error) {
	for i := 2; i >= 0; i-- {
		db, err = mysql.Open(ctx, dsn)
		if err != nil {
			logger.Error("failed to connect to mysql", zap.Error(err), zap.Int("attempts_remaining", i))
			time.Sleep(15 * time.Second)
		}
	}
	return
}

func newGraph(ctx context.Context, logger *zap.Logger) *ent.Client {
	var mysqlDSN string
	if mysqlDSN = os.Getenv("PG_MYSQL_DSN"); mysqlDSN == "" {
		panic(fmt.Errorf("failed to connect to mysql: missing PG_MYSQL_DSN"))
	}

	maxConnections := 10
	if limit := os.Getenv("PG_MYSQL_CONN_LIMIT"); limit != "" {
		connLimit, err := strconv.Atoi(limit)
		if err != nil {
			logger.Error("Invalid value set for PG_MYSQL_CONN_LIMIT, using default", zap.Error(err), zap.String("PG_MYSQL_CONN_LIMIT", limit), zap.Int("default_limit", maxConnections))
		} else {
			maxConnections = connLimit
		}
	}
	maxIdle := 10
	if limit := os.Getenv("PG_MYSQL_IDLE_LIMIT"); limit != "" {
		idleLimit, err := strconv.Atoi(limit)
		if err != nil {
			logger.Error("Invalid value set for PG_MYSQL_IDLE_LIMIT, using default", zap.Error(err), zap.String("PG_MYSQL_IDLE_LIMIT", limit), zap.Int("default_limit", maxIdle))
		} else {
			maxIdle = idleLimit
		}
	}
	logger = logger.With(zap.Int("mysql_max_conns", maxConnections), zap.Int("mysql_max_idle", maxIdle))

	db, err := connect(ctx, logger, mysqlDSN)
	if err != nil {
		panic(fmt.Errorf("failed to connect to mysql: %w", err))
	}

	db.SetMaxIdleConns(maxIdle)
	db.SetMaxOpenConns(maxConnections)
	db.SetConnMaxLifetime(time.Hour)
	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB("mysql", db)

	client := ent.NewClient(ent.Driver(drv))
	if err := client.Schema.Create(ctx, migrate.WithGlobalUniqueID(true)); err != nil {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close MySQL connection", zap.Error(err))
		}
		panic(err)
	}
	logger.Info("Connected to MySQL")

	return client
}
