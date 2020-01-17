package ssh

import (
	"io"
	"os"

	"github.com/kcarretto/paragon/pkg/script"
	"github.com/pkg/sftp"
)

// SendToTarget uploads a file to the remote system using SFTP over SSH.
func (env Environment) SendToTarget(parser script.ArgParser) (script.Retval, error) {
	// TODO: Handle timeout
	// TODO: Configure write permissions
	// TODO: Configure ownership
	name, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}

	remotePath, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}

	content, err := env.Download(name)
	if err != nil {
		return nil, err
	}

	client, err := env.Remote.Connect()
	if err != nil {
		return nil, err
	}

	session, err := sftp.NewClient(client)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	f, err := session.OpenFile(remotePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return nil, err
	}

	return io.Copy(f, content)
}
