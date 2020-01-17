package ssh

import (
	"github.com/kcarretto/paragon/pkg/script"
)

// Exec runs a command on the remote system using an underlying ssh session.
func (env Environment) Exec(parser script.ArgParser) (script.Retval, error) {
	cmd, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}

	client, err := env.Remote.Connect()
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	// TODO: Handle timeout
	result, err := session.CombinedOutput(cmd)
	return string(result), err
}
