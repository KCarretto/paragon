// +build !dev,!debug

package main

import "go.uber.org/zap"

func newLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		logger = zap.NewNop()
	}
	return logger
}
