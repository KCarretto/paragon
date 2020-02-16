// Package assets provides functionality to load asset files that were packed into the binary during
// compilation.
package assets

import (
	"net/http"

	"github.com/kcarretto/paragon/pkg/script"
)

// Environment used to configure the behaviour of calls to the ssh library.
type Environment struct {
	Assets http.FileSystem
}

// Library prepares a new assets library for use within a script environment.
func (env *Environment) Library(options ...func(*Environment)) script.Library {
	for _, opt := range options {
		opt(env)
	}

	return script.Library{
		"openFile": script.Func(env.openFile),
		// "exec":     script.Func(env.exec),
		// "openFile": script.Func(env.openFile),
	}
}

// Include the assets library in a script environment.
func (env *Environment) Include(options ...func(*Environment)) script.Option {
	return script.WithLibrary("assets", (*Environment).Library(env, options...))
}

// package assets

// import (
// 	"net/http"

// 	"github.com/kcarretto/paragon/pkg/script"
// )

// // Import the assets library to enable scripts to load assets from the provided filesystem.
// func Import(assets http.FileSystem) script.Library {
// 	return script.Library{
// 		"load": Load(assets),
// 	}
// }
