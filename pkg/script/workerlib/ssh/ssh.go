// Package ssh provides functionality to execute commands on remote targets using SSH. The library
// also provides standardized file upload and download methods that will rely on the SFTP protocol.
package ssh

import (
	"fmt"
	"log"

	"github.com/kcarretto/paragon/pkg/cdn"
	"github.com/kcarretto/paragon/pkg/script"
	"golang.org/x/crypto/ssh"
)

// Import the ssh library to enable scripts to access remote targets using ssh.
func Import(env Environment) script.Library {
	return script.Library{
		"exec":           script.Func(env.Exec),
		"sendToTarget":   script.Func(env.SendToTarget),
		"recvFromTarget": script.Func(env.RecvFromTarget),
	}
}

// Connector is responsible for establishing ssh sessions or reusing existing ones.
type Connector struct {
	Configs map[string][]*ssh.ClientConfig
}

func (conn *Connector) SetConfigs(remoteHost string, configs ...*ssh.ClientConfig) {
	if conn.Configs == nil {
		conn.Configs = make(map[string][]*ssh.ClientConfig)
	}

	conn.Configs[remoteHost] = configs
}

func (conn *Connector) Connect(host string) (*ssh.Client, error) {
	configs, ok := conn.Configs[host]
	if !ok || configs == nil {
		return nil, fmt.Errorf("no client configs found for host: %s", host)
	}

	// Append host
	host = host + ":22"

	// TODO: Implement caching
	for _, config := range configs {
		log.Printf("[DBG] Connecting to %s with config %+v", host, config)
		client, err := ssh.Dial("tcp", host, config)
		if err != nil {
			log.Printf("[ERR] Failed to %s: %+v", host, err)
			continue
		}

		return client, nil
	}

	return nil, fmt.Errorf("All client configs failed to connect")
}

// An Environment used for the execution of a single task.
type Environment struct {
	RemoteHost string
	*Connector
	cdn.Uploader
	cdn.Downloader
}
