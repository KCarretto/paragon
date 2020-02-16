package process

import (
	"fmt"

	"github.com/kcarretto/paragon/pkg/script"

	"go.starlark.net/starlark"
)

// Process provides a concurrency-safe wrapper for a Process that implements a starlark.Value.
type Process struct {
	Pid int32

	PPid    int32
	Name    string
	User    string
	Status  string
	CmdLine string
	Exe     string
	TTY     string
}

/*
 * Starlark.Value methods
 */

// String returns process metadata.
func (p Process) String() string {
	return fmt.Sprintf("%d %s %s", p.Pid, p.Name, p.User)
}

// Freeze is a no-op since the underlying process is safe for concurrent use.
func (p Process) Freeze() {}

// Type returns 'process' to indicate the type of the process within starlark.
func (p Process) Type() string {
	return "process"
}

// Truth value of a process is True if the PID is non-negative
func (p Process) Truth() starlark.Bool {
	if p.Pid < 0 {
		return starlark.False
	}
	return starlark.True
}

// Hash will error since the process type is not intended to be hashable.
func (p Process) Hash() (uint32, error) {
	return 0, fmt.Errorf("process type is unhashable")
}

// ParseParam from starlark input
func ParseParam(parser script.ArgParser, index int) (Process, error) {
	val, err := parser.GetParam(0)
	if err != nil {
		return Process{Pid: -1}, err
	}

	f, ok := val.(Process)
	if !ok {
		return Process{Pid: -1}, fmt.Errorf("%w: expected process type", script.ErrInvalidArgType)
	}

	return f, nil
}
