package env

import (
	"github.com/kcarretto/paragon/pkg/script"
)

// Environment TODO
type Environment struct {
	*script.Environment

	PrimaryIP       string
	OperatingSystem string
}

// Library TODO
func (env *Environment) Library(options ...func(*Environment)) script.Library {
	if env == nil {
		env = &Environment{}
	}

	for _, opt := range options {
		opt(env)
	}

	if env.Environment == nil {
		env.Environment = &script.Environment{}
	}

	return script.Library{
		"OS":        script.Func(env.os),
		"IP":        script.Func(env.ip),
		"isLinux":   script.Func(env.isLinux),
		"isWindows": script.Func(env.isWindows),
	}
}

// Include TODO
func (env *Environment) Include(options ...func(*Environment)) script.Option {
	return script.WithLibrary("target", (*Environment).Library(env, options...))
}
