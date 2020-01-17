// Package ssh provides functionality to execute commands on remote targets using SSH. The library
// also provides standardized file upload and download methods that will rely on the SFTP protocol.
package ssh

import (
	"io"

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

// A Connector is responsible for establishing ssh sessions or reusing existing ones.
type Connector interface {
	Connect() (*ssh.Client, error)
}

// An Uploader uploads files that were received from the remote system. It is responsible for
// closing the file after writing if the file implements io.Closer.
type Uploader interface {
	Upload(name string, file io.Reader) error
}

// A Downloader provides a file by name that will be sent to the remote system.
type Downloader interface {
	Download(name string) (io.Reader, error)
}

// An Environment used for the execution of a single task.
type Environment struct {
	Remote Connector
	Uploader
	Downloader
}
