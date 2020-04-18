// Package assets provides functionality to load asset files that were packed into the binary during
// compilation.
package assets

import (
	"github.com/kcarretto/paragon/pkg/cdn"
	"github.com/kcarretto/paragon/pkg/script"
	"github.com/spf13/afero"
)

// Environment used to configure the behaviour of calls to the ssh library.
type Environment struct {
	Assets     afero.Fs
	Files      []NamedReader
	Downloader cdn.Downloader
}

// Library prepares a new assets library for use within a script environment.
func (env *Environment) Library(options ...func(*Environment)) script.Library {
	for _, opt := range options {
		opt(env)
	}

	return script.Library{
		"file":    script.Func(env.file),
		"require": script.Func(env.require),
	}
}

// Include the assets library in a script environment.
func (env *Environment) Include(options ...func(*Environment)) script.Option {
	return script.WithLibrary("assets", (*Environment).Library(env, options...))
}
