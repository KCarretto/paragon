// Package ssh provides functionality to execute commands on remote targets using SSH. The library
// also provides standardized file upload and download methods that will rely on the SFTP protocol.
package ssh

import (
	"io"

	"github.com/kcarretto/paragon/pkg/script"
	"golang.org/x/crypto/ssh"
)

// Connector provides an ssh client connected to a remote host.
type Connector interface {
	Connect(remoteHost string, filter func([]ssh.ClientConfig) []ssh.ClientConfig) (*ssh.Client, error)
}

// Environment used to configure the behaviour of calls to the ssh library.
type Environment struct {
	*script.Environment

	Connector  Connector
	RemoteHost string
	RemoteUser string

	handles []io.Closer
}

// Library prepares a new ssh library for use within a script environment.
func (env *Environment) Library(options ...func(*Environment)) script.Library {
	if env == nil {
		env = &Environment{
			RemoteHost: "127.0.0.1:22",
		}
	}

	for _, opt := range options {
		opt(env)
	}

	if env.Environment == nil {
		env.Environment = &script.Environment{}
	}

	return script.Library{
		"setUser":       script.Func(env.setUser),
		"exec":          script.Func(env.exec),
		"openFile":      script.Func(env.openFile),
		"getRemoteHost": script.Func(env.getRemoteHost),
	}
}

// Include the ssh library in a script environment.
func (env *Environment) Include(options ...func(*Environment)) script.Option {
	return script.WithLibrary("ssh", (*Environment).Library(env, options...))
}
