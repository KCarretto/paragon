package script

import "io"

// An Environment enables script libraries to receive dependencies via dependency injection and
// register hooks for cleanup after execution.
type Environment struct {
	handles []io.Closer
}

// TrackHandle ensures that the provided handle is closed after script execution.
func (env *Environment) TrackHandle(handle io.Closer) {
	env.handles = append(env.handles, handle)
}

// Close all handles opened by the environment.
func (env *Environment) Close() (err error) {
	if env == nil || env.handles == nil {
		return nil
	}

	for _, handle := range env.handles {
		if handle == nil {
			continue
		}

		if closeErr := handle.Close(); closeErr != nil {
			err = closeErr
		}
	}
	return
}
