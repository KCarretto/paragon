package transport

import (
	"io"
	"sort"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Registrar interface {
	Add(name string, factory Factory, options ...Option)
	Update(name string, options ...Option) error
	Remove(name string) error
}

// A Registry tracks transports and associated metadata.
type Registry struct {
	SortBy func(t1, t2 Meta) bool

	active   map[string]io.WriteCloser
	metadata map[string]*Meta
}

// Add a transport factory to the registry and configure it's metadata.
func (reg Registry) Add(name string, factory Factory, options ...Option) {
	meta := &Meta{
		Name:    name,
		Factory: factory,
	}
	for _, opt := range options {
		opt(meta)
	}

	reg.metadata[name] = meta
}

// Update transport metadata by applying options to the transport metadata with the provided name.
func (reg Registry) Update(name string, options ...Option) error {
	transport, ok := reg.metadata[name]
	if !ok {
		return errors.New("transport not registered")
	}

	for _, opt := range options {
		opt(transport)
	}
	return nil
}

// Remove a transport from the registry, closing it if open. If the transport is not registered,
// remove is a no-op.
func (reg Registry) Remove(name string) error {
	transport, ok := reg.active[name]
	if !ok {
		return nil
	}
	defer delete(reg.active, name)
	return transport.Close()
}

// Get returns an active transport with the given name or initializes a new one if none is found.
func (reg Registry) Get(name string, logger *zap.Logger, writer Tasker) (io.WriteCloser, error) {
	transport, ok := reg.active[name]
	if ok && transport != nil {
		return transport, nil
	}

	meta, ok := reg.metadata[name]
	if !ok || meta == nil || meta.Factory == nil {
		// TODO: Log removal
		delete(reg.metadata, name)
		return nil, errors.New("transport not registered")
	}

	var err error
	if reg.active[name], err = meta.Factory(logger, writer); err != nil {
		return nil, errors.Wrap(err, "transport failed to initialize")
	}

	return reg.active[name], nil
}

// List transports, sorted by the provided SortBy method or by priority if no method was provided.
func (reg Registry) List() []Meta {
	if reg.SortBy == nil {
		reg.SortBy = func(t1, t2 Meta) bool {
			return t1.Priority < t2.Priority
		}
	}

	transports := make([]Meta, 0, len(reg.metadata))
	for name, meta := range reg.metadata {
		if meta != nil && meta.Factory != nil {
			transports = append(transports, *meta)
		} else {
			reg.Remove(name)
		}
	}

	sort.Slice(transports, func(i, j int) bool {
		return reg.SortBy(transports[i], transports[j])
	})

	return transports
}
