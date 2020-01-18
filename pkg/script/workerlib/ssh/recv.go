package ssh

import (
	"github.com/kcarretto/paragon/pkg/script"
	"github.com/pkg/sftp"
)

// RecvFromTarget downloads a file from the remote system using SFTP over SSH.
func (env Environment) RecvFromTarget(parser script.ArgParser) (script.Retval, error) {
	// TODO: Handle timeout
	name, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}

	remotePath, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}

	client, err := env.Connect(env.RemoteHost)
	if err != nil {
		return nil, err
	}

	session, err := sftp.NewClient(client)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	f, err := session.Open(remotePath)
	if err != nil {
		return nil, err
	}

	if err := env.Upload(name, f); err != nil {
		return nil, err
	}

	return nil, nil
}
