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

func Run(logger *zap.Logger) {
	a := agent.New(
		logger,
		Executor("Stuff"),
	)
	a.Transports.Add("local", transport.FactoryFn(local.New))
	a.Run()
	defer func() {
		if err := a.Close(); err != nil {
			panic(err)
		}
	}()

}
