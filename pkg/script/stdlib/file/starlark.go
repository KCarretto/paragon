package file

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/kcarretto/paragon/pkg/script"
	"go.starlark.net/starlark"
)

// New constructs a new file type for use within a starlark environment.
func New(f File) Type {
	return Type{
		f,
		&sync.RWMutex{},
	}
}

// Type provides a concurrency-safe wrapper for a File that implements a starlark.Value.
type Type struct {
	File

	mu *sync.RWMutex
}

// Name provides a concurrency-safe method to interact with the underlying File method.
func (f Type) Name() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.File.Name()
}

// Chmod provides a concurrency-safe method to interact with the underlying File method.
func (f Type) Chmod(mode os.FileMode) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.File.Chmod(mode)
}

// Chown provides a concurrency-safe method to interact with the underlying File method.
func (f Type) Chown(uid, gid int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.File.Chown(uid, gid)
}

// Stat provides a concurrency-safe method to interact with the underlying File method.
func (f Type) Stat() (os.FileInfo, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.File.Stat()
}

// Read provides a concurrency-safe method to interact with the underlying File method.
func (f Type) Read(p []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.File.Read(p)
}

// Write provides a concurrency-safe method to interact with the underlying File method.
func (f Type) Write(p []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.File.Write(p)
}

// Move provides a concurrency-safe method to interact with the underlying File method.
func (f Type) Move(dstPath string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.File.Move(dstPath)
}

// Close provides a concurrency-safe method to close the underlying File.
func (f Type) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if closer, ok := f.File.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

/*
 * Starlark.Value methods
 */

// String returns the name of the file.
func (f Type) String() string {
	return f.Name()
}

// Freeze is a no-op since the underlying file is safe for concurrent use.
func (f Type) Freeze() {}

// Type returns 'file' to indicate the type of the file within starlark.
func (f Type) Type() string {
	return "file"
}

// Truth value of a file is True if the file is non-nil
func (f Type) Truth() starlark.Bool {
	if f.File == nil {
		return starlark.False
	}
	return starlark.True
}

// Hash will error since the file type is not intended to be hashable.
func (f Type) Hash() (uint32, error) {
	return 0, fmt.Errorf("file type is unhashable")
}

// parseFileParam from starlark input
func parseFileParam(parser script.ArgParser, index int) (Type, error) {
	val, err := parser.GetParam(index)
	if err != nil {
		return Type{}, err
	}

	f, ok := val.(Type)
	if !ok {
		return Type{}, fmt.Errorf("%w: expected file type", script.ErrInvalidArgType)
	}

	return f, nil
}
