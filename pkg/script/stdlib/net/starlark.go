package net

import (
	"fmt"
	"syscall"

	"github.com/kcarretto/paragon/pkg/script"

	"go.starlark.net/starlark"
)

// Connection provides a concurrency-safe wrapper for a Connection that implements a starlark.Value.
type Connection struct {
	Fd         uint32
	ConnFamily uint32
	ConnType   uint32
	Pid        int32
	LocalIP    string
	LocalPort  uint32
	RemoteIP   string
	RemotePort uint32
	Status     string
}

func (c Connection) getConnectionType() string {
	switch c.ConnFamily {
	case syscall.AF_UNIX:
		return "unix"
	case syscall.AF_INET:
		if c.ConnType == syscall.SOCK_STREAM {
			return "tcp4"
		}
		if c.ConnType == syscall.SOCK_DGRAM {
			return "udp4"
		}
		break
	case syscall.AF_INET6:
		if c.ConnType == syscall.SOCK_STREAM {
			return "tcp6"
		}
		if c.ConnType == syscall.SOCK_DGRAM {
			return "udp6"
		}
		break
	default:
		return "unknown"
	}
	return "unknown"
}

/*
 * Starlark.Value methods
 */

// String returns connection metadata.
func (c Connection) String() string {
	return fmt.Sprintf("PID: %d, Type: %s, Status: %s, Local:%s:%d, Remote: %s:%d",
		c.Pid,
		c.getConnectionType(),
		c.Status,
		c.LocalIP,
		c.LocalPort,
		c.RemoteIP,
		c.RemotePort,
	)
}

// Freeze is a no-op since the underlying connection is safe for concurrent use.
func (c Connection) Freeze() {}

// Type returns 'connection' to indicate the type of the connection within starlark.
func (c Connection) Type() string {
	return "connection"
}

// Truth value of a connection is True if the Status is not the empty string.
func (c Connection) Truth() starlark.Bool {
	if c.Status == "" {
		return starlark.False
	}
	return starlark.True
}

// Hash will error since the connection type is not intended to be hashable.
func (c Connection) Hash() (uint32, error) {
	return 0, fmt.Errorf("connection type is unhashable")
}

// parseConnectionParam from starlark input
func parseConnectionParam(parser script.ArgParser, index int) (Connection, error) {
	val, err := parser.GetParam(0)
	if err != nil {
		return Connection{}, err
	}

	f, ok := val.(Connection)
	if !ok {
		return Connection{}, fmt.Errorf("%w: expected connection type", script.ErrInvalidArgType)
	}

	return f, nil
}
