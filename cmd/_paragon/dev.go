// +build dev

package main

import (
	"time"

	"github.com/kcarretto/paragon/agent"
	"github.com/kcarretto/paragon/agent/transport"
	"github.com/kcarretto/paragon/agent/transport/http"
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
	bot.Transports.Add(
		"http",
		transport.FactoryFn(http.New),
		transport.SetInterval(time.Second*5),
		transport.SetJitter(time.Second*1),
	)
}
