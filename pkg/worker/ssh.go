package worker

import (
	"fmt"
	"log"

	"github.com/kcarretto/paragon/ent"

	"golang.org/x/crypto/ssh"
)

// SSHConnector is used to establish ssh connections for the worker's task execution environment.
type SSHConnector struct {
	Credentials []*ent.Credential
	client      *ssh.Client
}

func (conn *SSHConnector) Connect(host string, filter func([]ssh.ClientConfig) []ssh.ClientConfig) (*ssh.Client, error) {
	if conn.client != nil {
		return conn.client, nil
	}

	// TODO: Don't set default 22 unless no port is specified
	host = host + ":22"

	if conn.Credentials == nil {
		return nil, fmt.Errorf("no credentials available for host: %s", host)
	}

	configs := filter(conn.configs())
	if configs == nil {
		return nil, fmt.Errorf("no valid client configs could be created for host: %s", host)
	}
	for _, config := range configs {
		log.Printf("[DBG] Connecting to %s with config %+v", host, config)
		client, err := ssh.Dial("tcp", host, &config)
		if err != nil {
			log.Printf("[ERR] Failed to %s: %+v", host, err)
			continue
		}
		conn.client = client
		return client, nil
	}

	return nil, fmt.Errorf("all connection attempts failed for host: %s", host)
}

func (conn *SSHConnector) Close() error {
	if conn.client == nil {
		return nil
	}
	return conn.client.Close()
}

func (conn *SSHConnector) configs() (configs []ssh.ClientConfig) {
	// Build a map of user -> ssh.AuthMethod
	creds := make(map[string][]ssh.AuthMethod)
	for _, credential := range conn.Credentials {
		// Skip empty credentials
		if credential == nil {
			continue
		}

		// Attempt to create an ssh.AuthMethod from the credential.
		method := conn.getAuthMethod(credential)
		if method == nil {
			continue
		}

		// Upsert a list of auth methods for the user
		userCreds, ok := creds[credential.Principal]
		if !ok || userCreds == nil {
			userCreds = []ssh.AuthMethod{}
		}

		// Append the auth method to the user's list
		creds[credential.Principal] = append(userCreds, method)
	}

	// Build an SSH config for each user
	for user, authMethods := range creds {
		if authMethods == nil {
			continue
		}

		config := ssh.ClientConfig{
			User:            user,
			Auth:            authMethods,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		configs = append(configs, config)
	}

	return
}

func (conn *SSHConnector) getAuthMethod(credential *ent.Credential) ssh.AuthMethod {
	switch credential.Kind {
	case "password":
		return ssh.Password(credential.Secret)
	case "key":
		signer, err := ssh.ParsePrivateKey([]byte(credential.Secret))
		if err != nil {
			log.Printf("[ERR] Failed to parse SSH private key (id=%d): %s",
				credential.ID,
				err.Error(),
			)
			return nil
		}
		return ssh.PublicKeys(signer)
	}

	log.Printf("[WARN] Unable to use credential for ssh connector (id=%d, kind=%s)",
		credential.ID,
		credential.Kind,
	)
	return nil
}
