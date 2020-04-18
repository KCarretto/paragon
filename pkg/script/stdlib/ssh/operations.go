package ssh

import (
	"fmt"
	"log"

	"github.com/kcarretto/paragon/pkg/script"
	libfile "github.com/kcarretto/paragon/pkg/script/stdlib/file"

	"github.com/pkg/sftp"
	"github.com/spf13/afero/sftpfs"
	"golang.org/x/crypto/ssh"
)

func (env *Environment) environmentFilter(configs []ssh.ClientConfig) []ssh.ClientConfig {
	if env.RemoteUser == "" {
		return configs
	}
	var newConfig []ssh.ClientConfig
	for _, config := range configs {
		if config.User == env.RemoteUser {
			newConfig = append(newConfig, config)
		}
	}
	return newConfig
}

// SetUser sets the RemoteUser attribute to be used in the outgoing SSH Connection. WARNING: MUST BE
// CALLED BEFORE OTHER SSH CALLS TO WORK.
//
//go:generate go run ../gendoc.go -lib ssh -func setUser -param user@String -doc "SetUser sets the RemoteUser attribute to be used in the outgoing SSH Connection. WARNING: MUST BE CALLED BEFORE OTHER SSH CALLS TO WORK."
//
// @callable: 	ssh.setUser
// @param: 		user 	@string
//
// @usage: 		ssh.setUser("root")
func (env *Environment) SetUser(user string) {
	env.RemoteUser = user
}
func (env *Environment) setUser(parser script.ArgParser) (script.Retval, error) {
	user, err := parser.GetString(0)
	if err != nil {
		log.Printf("[Err] SetUser param error: %v", err)
		return nil, err
	}
	env.SetUser(user)
	return nil, nil
}

// Exec a command on the remote system using an underlying ssh session.
//
//go:generate go run ../gendoc.go -lib ssh -func exec -param cmd@String -param disown@?Bool -retval output@String -retval err@Error -doc "Exec a command on the remote system using an underlying ssh session."
//
// @callable: 	ssh.exec
// @param: 		cmd 	@string
// @param: 		disown 	@?string
// @retval:		output 	@string
// @retval:		err 	@Error
//
// @usage: 		output, err = ssh.exec("/bin/bash -c 'ls -al'")
func (env *Environment) Exec(cmd string, disown bool) (string, error) {
	if env.Connector == nil {
		return "", fmt.Errorf("environment has no SSH connector")
	}

	client, err := env.Connector.Connect(env.RemoteHost, env.environmentFilter)
	if err != nil {
		log.Printf("[Err] Failed to connect to remote host: %v", err)
		return "", err
	}

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	env.TrackHandle(session)

	if disown {
		return "", session.Start(cmd)
	}

	// TODO: Handle exec timeout
	result, err := session.CombinedOutput(cmd)
	return string(result), err
}

func (env *Environment) exec(parser script.ArgParser) (script.Retval, error) {
	log.Printf("[DBG] Executing command on remote host via ssh")

	err := parser.RestrictKwargs("disown")
	if err != nil {
		return nil, err
	}

	cmd, err := parser.GetString(0)
	if err != nil {
		log.Printf("[Err] Exec param error: %v", err)
		return nil, err
	}

	disown, _ := parser.GetBoolByName("disown")

	retVal, retErr := env.Exec(cmd, disown)
	return script.WithError(retVal, retErr), nil
}

//go:generate go run ../gendoc.go -lib ssh -func file -param path@String -retval f@File -retval err@Error -doc "Prepare a descriptor for a file on the remote system using SFTP via SSH. The descriptor may be used with the file library."
func (env *Environment) file(parser script.ArgParser) (script.Retval, error) {
	path, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}

	client, err := env.Connector.Connect(env.RemoteHost, env.environmentFilter)
	if err != nil {
		return script.WithError(nil, err), nil
	}

	session, err := sftp.NewClient(client)
	if err != nil {
		return script.WithError(nil, err), nil
	}
	env.TrackHandle(session)

	return script.WithError(
		&libfile.File{
			Path: path,
			Fs:   sftpfs.New(session),
		},
		nil,
	), nil
}

// GetRemoteHost will return the remote host being used by the worker to connect to.
//
//go:generate go run ../gendoc.go -lib ssh -func getRemoteHost -retval host@String -doc "GetRemoteHost will return the remote host being used by the worker to connect to."
//
// @callable: 	ssh.getRemoteHost
// @retval:		host 		@String
//
// @usage: 		host = ssh.getRemoteHost()
func (env *Environment) GetRemoteHost() string {
	return env.RemoteHost
}

func (env *Environment) getRemoteHost(parser script.ArgParser) (script.Retval, error) {
	return env.GetRemoteHost(), nil
}
