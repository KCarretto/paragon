package websockets

import (
	"github.com/kcarretto/paragon/pkg/script"
)

var debug bool = true
var timeoutMS int = 2000
var parallelism int = 1000
var portSelection string

var hideUnavailableHosts bool
var versionRequested bool

var uuid = "30064771073"

func giveshell(parser script.ArgParser) (script.Retval, error) {
	websocket_host, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	shell_cmd, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}
	websocket_path, err := parser.GetString(2)
	if err != nil {
		websocket_path = "/cmd"
		// return nil, err
	}
	websocket_scheme, err := parser.GetString(3)
	if err != nil {
		websocket_scheme = "ws"
		// return nil, err
	}

	retVal, retErr := Giveshell(websocket_host, shell_cmd, websocket_path, websocket_scheme)
	return script.WithError(retVal, retErr), nil
}

//websockets.giveshell("127.0.0.1:4444", "bash", "-i", "")
func Giveshell(websocket_host string, shell_cmd string, shell_args string, websocket_scheme string) (string, error) {
	return "okay", nil
}

