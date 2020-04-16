package env

import (
	"math/rand"
	os_user "os/user"
	"runtime"
	"strings"
	"time"

	"github.com/kcarretto/paragon/pkg/script"
)

// UID returns the current user id. If not found, an empty string is returned.
//
//go:generate go run ../gendoc.go -lib env -retval uid@string -func uid -doc "UID returns the current user id. If not found, an empty string is returned."
//
// @callable: 	env.uid()
//
// @usage: 		uid = env.uid()
func (env *Environment) UID() string {
	usr, err := os_user.Current()
	if err != nil {
		return ""
	}

	return usr.Uid
}
func (env *Environment) uid(parser script.ArgParser) (script.Retval, error) {
	return env.UID(), nil
}

// User returns the current username. If not found, an empty string is returned.
//
//go:generate go run ../gendoc.go -lib env -retval username@string -func user -doc "User returns the current username. If not found, an empty string is returned."
//
// @callable: 	env.user()
//
// @usage: 		username = env.user()
func (env *Environment) User() string {
	usr, err := os_user.Current()
	if err != nil {
		return ""
	}

	return usr.Username
}
func (env *Environment) user(parser script.ArgParser) (script.Retval, error) {
	return env.User(), nil
}

// Time returns the current number of seconds since the unix epoch.
//
//go:generate go run ../gendoc.go -lib env -retval i@int -func time -doc "Time returns the current number of seconds since the unix epoch."
//
// @callable: 	env.time()
//
// @usage: 		now = env.time()
func (env *Environment) Time() int {
	return int(time.Now().Unix())
}
func (env *Environment) time(parser script.ArgParser) (script.Retval, error) {
	return env.Time(), nil
}

// Rand returns a random int. Not cryptographically secure.
//
//go:generate go run ../gendoc.go -lib env -retval i@int -func rand -doc "Rand returns a random int. Not cryptographically secure."
//
// @callable: 	env.rand()
//
// @usage: 		num = env.rand()
func (env *Environment) Rand() int {
	return rand.Int()
}
func (env *Environment) rand(parser script.ArgParser) (script.Retval, error) {
	return env.Rand(), nil
}

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
		env.OperatingSystem = strings.ToUpper(runtime.GOOS)
	}

	return strings.ToUpper(env.OperatingSystem)
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
		env.OperatingSystem = strings.ToUpper(runtime.GOOS)
	}

	return strings.ToUpper(env.OperatingSystem) == "LINUX"
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
		env.OperatingSystem = strings.ToUpper(runtime.GOOS)
	}

	return strings.ToUpper(env.OperatingSystem) == "WINDOWS"
}
func (env *Environment) isWindows(parser script.ArgParser) (script.Retval, error) {
	return env.IsWindows(), nil
}
