package transport

import (
	"errors"
	"io"
	"sort"
	"sync"
)

// A Registry tracks transports and associated metadata.
type Registry struct {
	SortBy func(t1, t2 Transport) bool

	mu         sync.RWMutex
	transports map[string]Transport
}

func (reg *Registry) set(transport Transport) {
	reg.mu.Lock()
	if reg.transports == nil {
		reg.transports = map[string]Transport{}
	}
	reg.transports[transport.Name] = transport
	reg.mu.Unlock()
}

// Add a transport to the registry. If a transport with the same name is already registered, this
// method is a no-op. Use Update() to modify existing transports.
func (reg *Registry) Add(transport Transport) {
	if _, ok := reg.transports[transport.Name]; ok {
		return
	}

	reg.set(transport)
}

// Get returns a registered transport with the given name.
func (reg *Registry) Get(name string) (Transport, error) {
	reg.mu.RLock()
	transport, ok := reg.transports[name]
	reg.mu.RUnlock()

	if !ok {
		return Transport{}, errors.New("transport not registered")
	}

	return transport, nil
}

// Update the transport with the provided name by applying options to it.
func (reg *Registry) Update(name string, options ...Option) error {
	transport, err := reg.Get(name)
	if err != nil {
		return err
	}

	for _, opt := range options {
		opt(&transport)
	}

	reg.set(transport)

	return nil
}

// List transports, sorted by the provided SortBy method or by priority if no method was provided.
func (reg *Registry) List() []Transport {
	if reg.SortBy == nil {
		reg.SortBy = func(t1, t2 Transport) bool {
			return t1.Priority < t2.Priority
		}
	}

	transports := make([]Transport, 0, len(reg.transports))
	for _, transport := range reg.transports {
		transports = append(transports, transport)
	}

	sort.Slice(transports, func(i, j int) bool {
		return reg.SortBy(transports[i], transports[j])
	})

	return transports
}

// Close a transport by name. Close is a no-op if the transport is not an io.Closer.
func (reg *Registry) Close(name string) error {
	transport, err := reg.Get(name)
	if err != nil {
		return err
	}

	if closer, ok := transport.Writer.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}

// CloseAll attempts to close all transports.
func (reg *Registry) CloseAll() {
	for name := range reg.transports {
		reg.Close(name)
	}
}
