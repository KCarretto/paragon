package transport

import (
	"io"
	"sort"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

// A Factory is used to initialize a new transport.
type Factory func(*zap.Logger, TaskWriter) (io.WriteCloser, error)

// A TaskWriter is provided to a Transport to enable it to send incoming tasks to the agent.
type TaskWriter interface {
	WriteTask(id string, content io.Reader)
}

// An Option enables additional configuration of a transport Meta.
type Option func(*Meta)

// Meta holds metadata about a transport and a factory method to initialize it.
type Meta struct {
	Name     string
	Priority int
	Interval time.Duration
	Jitter   time.Duration
	Factory  Factory
}

// A Registry maintains a mapping of transport metadata registered with it and provides functionality
// to return a sorted list of registered transport metadata.
type Registry struct {
	SortBy func(t1, t2 Meta) bool

	active   map[string]io.WriteCloser
	metadata map[string]*Meta
}

// Register a transport factory and configure it's metadata.
func (reg Registry) Register(name string, factory Factory, options ...Option) {
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

// Get returns a named transport io.WriteCloser.
func (reg Registry) Get(name string, logger *zap.Logger, writer TaskWriter) (io.WriteCloser, error) {
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

// CloseTransport closes an open transport. If the transport does not exist or is not active, Close is a no-op.
func (reg Registry) CloseTransport(name string) error {
	transport, ok := reg.active[name]
	if !ok {
		return nil
	}
	defer delete(reg.active, name)
	return transport.Close()
}

// Close closes all active transports.
func (reg Registry) Close() error {
	var errs []error
	for name, transport := range reg.active {
		if err := transport.Close(); err != nil {
			errs = append(errs, errors.Wrapf(err, "failed to close transport %s", name))
		}
	}
	return multierr.Combine(errs...)
}

// List transports, sorted by the provided SortBy method or by priority if none was provided.
func (reg Registry) List() []Meta {
	if reg.SortBy == nil {
		reg.SortBy = func(t1, t2 Meta) bool {
			return t1.Priority < t2.Priority
		}
	}

	transports := make([]Meta, 0, len(reg.metadata))
	for name, meta := range reg.metadata {
		if meta == nil || meta.Factory == nil {
			// TODO: Log removal
			delete(reg.metadata, name)
		} else {
			transports = append(transports, *meta)
		}
	}

	sort.Slice(transports, func(i, j int) bool {
		return reg.SortBy(transports[i], transports[j])
	})

	return transports
}

// SetPriority metadata for the transport.
func SetPriority(priority int) Option {
	return func(meta *Meta) {
		meta.Priority = priority
	}
}

// SetInterval metadata for the transport.
func SetInterval(interval time.Duration) Option {
	return func(meta *Meta) {
		meta.Interval = interval
	}
}

// SetJitter metadata for the transport.
func SetJitter(jitter time.Duration) Option {
	return func(meta *Meta) {
		meta.Jitter = jitter
	}
}
