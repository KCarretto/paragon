// Package stdlib contains all of the basic building blocks of OS operation. It is intended to be cross platform and
// minimalistic in scale. Its largest sublibrary is "sys" and to use it simply use the syntax `load("sys", "myFunc")`
// where `myFunc` is one of the premade library functions.
package stdlib

import (
	"net/http"

	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/assets"
	"github.com/kcarretto/paragon/pkg/script/stdlib/sys"
)

// An Option enables configuration of stdlib dependencies.
type Option func(*Dependencies)

// Dependencies holds library dependencies that may be accessed by stdlib functions.
type Dependencies struct {
	Assets http.FileSystem
}

// Load is a script Option that loads the standard library into a script's execution environment
func Load(options ...Option) script.Option {
	deps := Dependencies{}
	for _, opt := range options {
		opt(&deps)
	}

	libs := map[string]script.Library{
		"sys":    sys.Import(),
		"assets": assets.Import(deps.Assets),
	}

	return script.WithLibraries(libs)
}

// WithAssets loading assets from the provided filesystem.
func WithAssets(fs http.FileSystem) Option {
	return func(deps *Dependencies) {
		deps.Assets = fs
	}
}
