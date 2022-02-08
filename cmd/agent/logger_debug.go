//go:build dev || debug
// +build dev debug

package main

import (
	"go.uber.org/zap"
)

func newLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger
}
