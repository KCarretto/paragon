package env

import (
	"math/rand"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/kcarretto/paragon/pkg/script"
)

//go:generate go run ../gendoc.go -lib env -func IP -doc "IP returns the primary IP address."
func (env *Environment) ip(parser script.ArgParser) (script.Retval, error) {
	return env.PrimaryIP, nil
}

//go:generate go run ../gendoc.go -lib env -func OS -doc "OS returns the operating system."
func (env *Environment) os(parser script.ArgParser) (script.Retval, error) {
	if env.OperatingSystem == "" {
		env.OperatingSystem = strings.ToUpper(runtime.GOOS)
	}

	return strings.ToUpper(env.OperatingSystem), nil
}

//go:generate go run ../gendoc.go -lib env -retval pid@Int -func PID -doc "PID returns the id of the current process."
func (env *Environment) pid(parser script.ArgParser) (script.Retval, error) {
	return os.Getpid(), nil
}

//go:generate go run ../gendoc.go -lib env -retval uid@String -func UID -doc "UID returns the current user id. If not found, an empty string is returned."
func (env *Environment) uid(parser script.ArgParser) (script.Retval, error) {
	usr, err := user.Current()
	if err != nil {
		return "", nil
	}

	return usr.Uid, nil
}

//go:generate go run ../gendoc.go -lib env -retval username@string -func user -doc "user returns the current username. If not found, an empty string is returned."
func (env *Environment) user(parser script.ArgParser) (script.Retval, error) {
	usr, err := user.Current()
	if err != nil {
		return "", nil
	}

	return usr.Username, nil
}

//go:generate go run ../gendoc.go -lib env -retval i@int -func time -doc "time returns the current number of seconds since the unix epoch."
func (env *Environment) time(parser script.ArgParser) (script.Retval, error) {
	return int(time.Now().Unix()), nil
}

//go:generate go run ../gendoc.go -lib env -retval i@int -func rand -doc "rand returns a random int. Not cryptographically secure."
func (env *Environment) rand(parser script.ArgParser) (script.Retval, error) {
	return rand.Int(), nil
}

//go:generate go run ../gendoc.go -lib env -func isLinux -doc "isLinux returns true if the operating system is linux."
func (env *Environment) isLinux(parser script.ArgParser) (script.Retval, error) {
	if env.OperatingSystem == "" {
		env.OperatingSystem = strings.ToUpper(runtime.GOOS)
	}

	return strings.ToUpper(env.OperatingSystem) == "LINUX", nil
}

//go:generate go run ../gendoc.go -lib env -func isWindows -doc "isWindows returns true if the operating system is windows."
func (env *Environment) isWindows(parser script.ArgParser) (script.Retval, error) {
	if env.OperatingSystem == "" {
		env.OperatingSystem = strings.ToUpper(runtime.GOOS)
	}

	return strings.ToUpper(env.OperatingSystem) == "WINDOWS", nil
}
