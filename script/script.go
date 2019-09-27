package script

import (
	"bytes"
	"io"
)

// A Script provides metadata and instructions to be executed by the interpreter.
type Script struct {
	io.Reader        // Read() instructions to execute
	ID        string // Describes the ID of the script
}

// New initializes and returns a script with the provided contents.
func New(id string, content []byte) *Script {
	return &Script{
		bytes.NewBuffer(content),
		id,
	}
}
