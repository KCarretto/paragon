// +build dev

package main

import (
	"bytes"
	"time"

	"github.com/kcarretto/paragon/script"

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

func Run(logger *zap.Logger) {
	py := script.NewInterpreter()
	a := agent.New(
		logger,
		Executor{
			py,
		},
	)
	a.Transports.Add(
		"http",
		transport.FactoryFn(http.New),
		transport.SetInterval(time.Second*5),
		transport.SetJitter(time.Second*1),
	)

	a.Run()
	a.QueueTask("someID", bytes.NewBufferString(`print("Hello World")`))
	defer func() {
		if err := a.Close(); err != nil {
			panic(err)
		}
	}()

}
