package transport

import (
	"io"
	"sort"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// A Registry tracks transports and associated metadata.
type Registry struct {
	SortBy func(t1, t2 Meta) bool

	mu       sync.RWMutex
	active   map[string]io.WriteCloser
	metadata map[string]Meta
}

// Add a transport factory to the registry and configure it's metadata. If a transport with the same
// name is already registered, this method is a no-op. Use Update() and Remove() to modify existing
// transports.
func (reg Registry) Add(name string, factory Factory, options ...Option) {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	if _, ok := reg.metadata[name]; ok {
		return
	}

	meta := Meta{
		Name:    name,
		Factory: factory,
	}
	for _, opt := range options {
		opt(&meta)
	}

	if reg.metadata == nil {
		reg.metadata = map[string]Meta{}
	}
	reg.metadata[name] = meta
}

// Update transport metadata by applying options to the transport metadata with the provided name.
func (reg Registry) Update(name string, options ...Option) error {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	meta, ok := reg.metadata[name]
	if !ok {
		return errors.New("transport not registered")
	}

	for _, opt := range options {
		opt(&meta)
	}
	return nil
}

// Remove a transport from the registry, closing it if open. If the transport is not registered,
// remove is a no-op.
func (reg Registry) Remove(name string) error {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	transport, ok := reg.active[name]
	if !ok {
		return nil
	}
	delete(reg.active, name)

	return transport.Close()
}

// Get returns an active transport with the given name or initializes a new one if none is found.
func (reg Registry) Get(name string, logger *zap.Logger, writer Tasker) (io.Writer, error) {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	transport, ok := reg.active[name]
	if ok && transport != nil {
		return transport, nil
	}

	meta, ok := reg.metadata[name]
	if !ok || meta.Factory == nil {
		delete(reg.metadata, name)
		return nil, errors.New("transport not registered")
	}

	transport, err := meta.Factory.New(logger, writer)
	if err != nil {
		return nil, errors.Wrap(err, "transport failed to initialize")
	}
	if transport == nil {
		err = errors.New("transport factory returned nil")
		logger.DPanic("Transport factory cannot return nil transport without error", zap.Error(err))
		return nil, err
	}

	if reg.active == nil {
		reg.active = map[string]io.WriteCloser{}
	}
	reg.active[name] = transport

	return transport, nil
}

// List transports, sorted by the provided SortBy method or by priority if no method was provided.
func (reg Registry) List() []Meta {
	reg.mu.RLock()
	defer reg.mu.RUnlock()

	if reg.SortBy == nil {
		reg.SortBy = func(t1, t2 Meta) bool {
			return t1.Priority < t2.Priority
		}
	}

	transports := make([]Meta, 0, len(reg.metadata))
	for name, meta := range reg.metadata {
		if meta.Factory != nil {
			transports = append(transports, meta)
		} else {
			reg.Remove(name)
		}
	}

	sort.Slice(transports, func(i, j int) bool {
		return reg.SortBy(transports[i], transports[j])
	})

	return transports
}
