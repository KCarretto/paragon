package ssh

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/file"
	"github.com/pkg/sftp"
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

// OpenFile on the remote system using SFTP over SSH. The file is created if it does not yet exist.
//
//go:generate go run ../gendoc.go -lib ssh -func openFile -param path@String -retval f@File -retval err@Error -doc "OpenFile on the remote system using SFTP over SSH. The file is created if it does not yet exist."
//
// @callable: 	ssh.openFile
// @param: 		path	 	@string
// @retval:		file 		@File
// @retval:		err 		@Error
//
// @usage: 		f, err = ssh.openFile("/bin/implant")
func (env *Environment) OpenFile(filePath string) (file.Type, error) {
	client, err := env.Connector.Connect(env.RemoteHost, env.environmentFilter)
	if err != nil {
		return file.Type{}, err
	}

	session, err := sftp.NewClient(client)
	if err != nil {
		return file.Type{}, err
	}
	env.TrackHandle(session)
	dir := path.Dir(filePath)
	if err := session.MkdirAll(dir); err != nil {
		return file.New(nil), fmt.Errorf("failed to create parent directory %q: %w", dir, err)
	}

	f, err := session.OpenFile(filePath, os.O_RDWR|os.O_CREATE)
	if err != nil {
		return file.Type{}, err
	}
	handle := &File{f, session}
	env.TrackHandle(handle)

	return file.New(handle), nil
}

func (env *Environment) openFile(parser script.ArgParser) (script.Retval, error) {
	filePath, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}

	retVal, retErr := env.OpenFile(filePath)
	return script.WithError(retVal, retErr), nil
}
