package main

import (
	"github.com/kcarretto/paragon/agent"
)

func NewAgent() *agent.Agent {
	return agent.New(
		getLogger(),
		getExecutor(),
	)
}
