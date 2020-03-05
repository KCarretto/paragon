package http

import (
	"fmt"

	"github.com/kcarretto/paragon/pkg/script"

	"go.starlark.net/starlark"
)

// Request provides a concurrency-safe wrapper for an HTTPRequest that implements a starlark.Value.
type Request struct {
	Url     string
	Method  string
	Headers map[string]string
	Body    string
}

/*
 * Starlark.Value methods
 */

// String returns Request metadata.
func (r Request) String() string {
	return fmt.Sprintf("%s %s", r.Url, r.Method)
}

// Freeze is a no-op since the underlying Request is safe for concurrent use.
func (r Request) Freeze() {}

// Type returns 'httprequest' to indicate the type of the Request within starlark.
func (r Request) Type() string {
	return "httprequest"
}

// Truth value of a Request is True if the PID is non-negative
func (r Request) Truth() starlark.Bool {
	if r.Url == "" {
		return starlark.False
	}
	return starlark.True
}

// Hash will error since the Request type is not intended to be hashable.
func (r Request) Hash() (uint32, error) {
	return 0, fmt.Errorf("httprequest type is unhashable")
}

// ParseParam from starlark input
func ParseParam(parser script.ArgParser, index int) (*Request, error) {
	val, err := parser.GetParam(0)
	if err != nil {
		return &Request{Url: ""}, err
	}

	r, ok := val.(*Request)
	if !ok {
		return &Request{Url: ""}, fmt.Errorf("%w: expected Request type", script.ErrInvalidArgType)
	}

	return r, nil
}
