package ssh

import (
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
)

// Connector is responsible for establishing ssh sessions or reusing existing ones.
type Connector struct {
	Configs map[string][]*ssh.ClientConfig
	clients []*ssh.Client
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
		conn.clients = append(conn.clients, client)
		return client, nil
	}

	return nil, fmt.Errorf("All client configs failed to connect")
}

// Close all clients opened by this connector.
func (conn *Connector) Close() (err error) {
	for _, client := range conn.clients {
		if closeErr := client.Close(); closeErr != nil {
			err = closeErr
		}
	}

	return err
}
