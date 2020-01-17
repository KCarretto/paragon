// Package workerlib encompasses the set of functions that a worker may execute, like uploading
// files to target machines, executing commands over SSH, WinRM, etc.
package workerlib

import (
	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/workerlib/ssh"
)

// An Option enables configuration of workerlib dependencies.
type Option func(*Dependencies)

// Dependencies holds library dependencies that may be accessed by workerlib functions.
type Dependencies struct {
	SSH ssh.Environment
}

// Load is a script Option that loads worker libraries into a script's execution environment.
func Load(options ...Option) script.Option {
	deps := Dependencies{}
	for _, opt := range options {
		opt(&deps)
	}

	libs := map[string]script.Library{
		"ssh": ssh.Import(deps.SSH),
	}

	return script.WithLibraries(libs)
}

// WithSSH configures the dependencies used by SSH methods.
func WithSSH(env ssh.Environment) Option {
	return func(deps *Dependencies) {
		deps.SSH = env
	}
}
