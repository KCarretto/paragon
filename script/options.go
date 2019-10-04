package script

import "io"

// An Option enables additional customization of script configuration and execution.
type Option func(*Script)

// WithOutput sets the destination for script execution output.
func WithOutput(w io.Writer) Option {
	return func(script *Script) {
		script.Writer = w
	}
}

// WithLibrary is an option to add a library to the script's execution environment.
func WithLibrary(name string, lib Library) Option {
	return func(script *Script) {
		if script.Libraries == nil {
			script.Libraries = map[string]Library{}
		}

		script.Libraries[name] = lib
	}
}

// WithLibraries adds one or more libraries to the script's execution environment.
func WithLibraries(libs map[string]Library) Option {
	var opts []Option
	for name, lib := range libs {
		opts = append(opts, WithLibrary(name, lib))
	}

	return func(script *Script) {
		for _, opt := range opts {
			opt(script)
		}
	}
}
