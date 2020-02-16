package ssh

import (
	"fmt"
	"log"
	"os"

	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/file"
	"github.com/pkg/sftp"
)

// Exec a command on the remote system using an underlying ssh session.
//
// @callable: 	ssh.exec
// @param: 		cmd 	@string
// @retval:		output 	@string
// @retval:		err 	@Error
//
// @usage: 		output, err = ssh.exec("/bin/bash -c 'ls -al'")
func (env *Environment) Exec(cmd string) (string, error) {
	if env.Connector == nil {
		return "", fmt.Errorf("environment has no SSH connector")
	}

	client, err := env.Connector.Connect(env.RemoteHost)
	if err != nil {
		log.Printf("[Err] Failed to connect to remote host: %v", err)
		return "", err
	}

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	env.TrackHandle(session)

	// TODO: Handle exec timeout
	result, err := session.CombinedOutput(cmd)
	return string(result), err
}

func (env *Environment) exec(parser script.ArgParser) (script.Retval, error) {
	log.Printf("[DBG] Executing command on remote host via ssh")

	cmd, err := parser.GetString(0)
	if err != nil {
		log.Printf("[Err] Exec param error: %v", err)
		return nil, err
	}

	retVal, retErr := env.Exec(cmd)
	return script.WithError(retVal, retErr), nil
}

// OpenFile on the remote system using SFTP over SSH. The file is created if it does not yet exist.
//
// @callable: 	ssh.openFile
// @param: 		path 	@string
// @retval:		file 	@File
// @retval:		err 	@Error
//
// @usage: 		f, err = ssh.openFile("/bin/implant")
func (env *Environment) OpenFile(path string) (file.Type, error) {
	client, err := env.Connector.Connect(env.RemoteHost)
	if err != nil {
		return file.Type{}, err
	}

	session, err := sftp.NewClient(client)
	if err != nil {
		return file.Type{}, err
	}
	env.TrackHandle(session)

	f, err := session.OpenFile(path, os.O_RDWR|os.O_CREATE)
	if err != nil {
		return file.Type{}, err
	}
	handle := &File{
		f,
		session,
	}
	env.TrackHandle(handle)

	return file.New(handle), nil
}

func (env *Environment) openFile(parser script.ArgParser) (script.Retval, error) {
	path, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}

	retVal, retErr := env.OpenFile(path)
	return script.WithError(retVal, retErr), nil
}
