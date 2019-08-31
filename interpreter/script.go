package interpreter

import (
	"bytes"
	"io"
)

// A Script provides metadata and instructions to be executed by the interpreter.
type Script struct {
	io.Reader        // Read() instructions to execute
	Name      string // Describes the name of the script
}

// NewScript initializes and returns a script with the provided contents.
func NewScript(name string, content []byte) *Script {
	return &Script{
		Reader: bytes.NewBuffer(content),
		Name:   name,
	}
}
