// +build local

package main

import (
	"github.com/kcarretto/paragon/agent"
	"github.com/kcarretto/paragon/agent/transport"
	"github.com/kcarretto/paragon/agent/transport/local"
	"go.uber.org/zap"
)

func getLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger
}

func addTransports(bot *agent.Agent) {
	bot.Transports.Add("local", transport.FactoryFn(local.New))
}
