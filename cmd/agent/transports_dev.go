// +build dev

package main

import (
	"github.com/kcarretto/paragon/agent"
	"go.uber.org/zap"
)

func transports(logger *zap.Logger) (transports []agent.Transport) {
	// TODO: Add localhost HTTP transport
	return
}
