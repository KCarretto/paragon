// Package ssh provides functionality to execute commands on remote targets using SSH. The library
// also provides standardized file upload and download methods that will rely on the SFTP protocol.
package ssh

import (
	"github.com/kcarretto/paragon/pkg/script"
)

// Environment used to configure the behaviour of calls to the ssh library.
type Environment struct {
	RemoteHost string
	*Connector
}

// Library prepares a new ssh library for use within a script environment.
func Library(options ...func(*Environment)) script.Library {
	env := &Environment{
		RemoteHost: "127.0.0.1:22",
		Connector:  &Connector{},
	}
	for _, opt := range options {
		opt(env)
	}

	return script.Library{
		"exec":     script.Func(env.exec),
		"openFile": script.Func(env.openFile),
	}
}

// Include the ssh library in a script environment.
func Include(options ...func(*Environment)) script.Option {
	return script.WithLibrary("ssh", Library(options...))
}
