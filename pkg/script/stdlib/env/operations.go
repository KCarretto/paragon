package env

import (
	"runtime"

	"github.com/kcarretto/paragon/pkg/script"
)

// IP returns the primary IP address.
//
//go:generate go run ../gendoc.go -lib env -func IP -doc "IP returns the primary IP address."
//
// @callable: 	env.IP()
//
// @usage: 		ip = env.IP()
func (env *Environment) IP() string {
	return env.PrimaryIP
}
func (env *Environment) ip(parser script.ArgParser) (script.Retval, error) {
	return env.IP(), nil
}

// OS returns the operating system.
//
//go:generate go run ../gendoc.go -lib env -func OS -doc "OS returns the operating system."
//
// @callable: 	env.OS()
//
// @usage: 		os = env.OS()
func (env *Environment) OS() string {
	if env.OperatingSystem == "" {
		env.OperatingSystem = runtime.GOOS
	}

	return env.OperatingSystem
}
func (env *Environment) os(parser script.ArgParser) (script.Retval, error) {
	return env.OS(), nil
}

// IsLinux returns true if the operating system is linux.
//
//go:generate go run ../gendoc.go -lib env -func isLinux -doc "IsLinux returns true if the operating system is linux."
//
// @callable: 	env.isLinux()
//
// @usage: 		env.isLinux()
func (env *Environment) IsLinux() bool {
	if env.OperatingSystem == "" {
		env.OperatingSystem = runtime.GOOS
	}

	return env.OperatingSystem == "linux"
}
func (env *Environment) isLinux(parser script.ArgParser) (script.Retval, error) {
	return env.IsLinux(), nil
}

// IsWindows returns true if the operating system is windows.
//
//go:generate go run ../gendoc.go -lib env -func isWindows -doc "IsWindows returns true if the operating system is windows."
//
// @callable: 	env.isWindows()
//
// @usage: 		env.isWindows()
func (env *Environment) IsWindows() bool {
	if env.OperatingSystem == "" {
		env.OperatingSystem = runtime.GOOS
	}

	return env.OperatingSystem == "windows"
}
func (env *Environment) isWindows(parser script.ArgParser) (script.Retval, error) {
	return env.IsWindows(), nil
}
