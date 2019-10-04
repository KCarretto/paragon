package sys

import (
	"runtime"

	"github.com/kcarretto/paragon/script"
)

// Detect uses runtime.GOOS to detect what OS the agent is running on.
//
// Enum for return value {
//  "linux"
//  "windows"
//  "darwin"
//  "bsd"
//  "solaris"
//  "other"
// }
//
// @return (osStr, nil) iff success; (nil, err) o/w
func Detect(parser script.ArgParser) (script.Retval, error) {
	var osStr string
	switch runtime.GOOS {
	case "linux":
		osStr = "linux"
	case "windows":
		osStr = "windows"
	case "darwin":
		osStr = "darwin"
	case "freebsd", "netbsd", "openbsd", "dragonfly":
		osStr = "bsd"
	case "android":
		osStr = "android"
	case "solaris":
		osStr = "solaris"
	default:
		osStr = "other"
	}
	return osStr, nil
}
