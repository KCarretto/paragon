package ssh

import (
	"log"

	"github.com/kcarretto/paragon/pkg/script"
)

// Exec runs a command on the remote system using an underlying ssh session.
func (env Environment) Exec(parser script.ArgParser) (script.Retval, error) {
	log.Printf("[DBG] Executing command on remote host via ssh")

	cmd, err := parser.GetString(0)
	if err != nil {
		log.Printf("[Err] Exec param error: %v", err)
		return nil, err
	}

	client, err := env.Connect(env.RemoteHost)
	if err != nil {
		log.Printf("[Err] Failed to connect to remote host: %v", err)
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
