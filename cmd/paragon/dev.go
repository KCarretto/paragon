// +build development

package main

import (
	"github.com/kcarretto/paragon/agent"
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
	agent := agent.New(
		logger,
		Executor("Stuff"),
	)

	agent.Run()
	defer func() {
		if err := agent.Close(); err != nil {
			panic(err)
		}
	}()

}
