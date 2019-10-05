package sys

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"syscall"

	"github.com/kcarretto/paragon/script"
	"github.com/shirou/gopsutil/net"
)

func familyAndTypeToString(connFamily, connType uint32) string {
	if connFamily == syscall.AF_UNIX {
		return "unix"
	}
	if connFamily == syscall.AF_INET {
		if connType == syscall.SOCK_STREAM {
			return "tcp4"
		}
		if connType == syscall.SOCK_DGRAM {
			return "udp4"
		}
	}
	if connFamily == syscall.AF_INET6 {
		if connType == syscall.SOCK_STREAM {
			return "tcp6"
		}
		if connType == syscall.SOCK_DGRAM {
			return "udp6"
		}
	}
	return "unknown"
}

// Connections uses net.ConnectionsPid to get all the connections opened by a process, given a connection protocol.
// May work on windows.
//
// The Connections map is defined as:
//  map[string]string{
//   "pid":           PID of the process,
//   "proto":         PID of the parent process,
//   "localaddr":     Name of the process,
//   "remoteaddr":    User who ran/owns the process,
//   "status":        Status of the process,
//  }
//
// @return (connections, nil) iff success; (nil, err) o/w
func Connections(parser script.ArgParser) (script.Retval, error) {
	err := parser.RestrictKwargs("type", "ppid")
	if err != nil {
		return nil, err
	}

	connectionsType, err := parser.GetStringByName("type")
	if err != nil {
		connectionsType = "all"
	}
	ppid, err := parser.GetIntByName("ppid")
	if err != nil {
		ppid = 0
	}
	connections, err := net.ConnectionsPid(connectionsType, int32(ppid))
	var conns []map[string]string
	for _, conn := range connections {
		connMap := map[string]string{
			"pid":        string(conn.Pid),
			"proto":      familyAndTypeToString(conn.Family, conn.Type),
			"localaddr":  conn.Laddr.IP + ":" + string(conn.Laddr.Port),
			"remoteaddr": conn.Raddr.IP + ":" + string(conn.Raddr.Port),
			"status":     conn.Status,
		}
		conns = append(conns, connMap)
	}
	return conns, nil
}

// Request uses the http package to send a HTTP request and return a response. Currently we only support GET and POST.
//
// @param requestURL:   A string that is the full requested URL.
// @param ?method:      The HTTP method of the request. The default is GET.
// @param ?writeToFile: The path to the file you wish to have the response written to. The default is to
// just return the string.
// @param ?contentType: The value of the "Content-Type" header for your request. Required for POST, the default is
// "application/json".
// @param ?data:        The body of your request. Only Available for POST.
//
// @return (response, nil) iff success; (nil, err) o/w
func Request(parser script.ArgParser) (script.Retval, error) {
	err := parser.RestrictKwargs("method", "writeToFile", "contentType", "data")
	if err != nil {
		return nil, err
	}

	requestURL, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	method, err := parser.GetStringByName("method")
	if errors.Is(err, script.ErrMissingKwarg) {
		method = "GET"
	} else if err != nil {
		fmt.Println(err)
		return nil, err
	}
	file, err := parser.GetStringByName("writeToFile")
	if errors.Is(err, script.ErrMissingKwarg) {
		file = ""
	} else if err != nil {
		return nil, err
	}
	var resp *http.Response
	switch method {
	case "POST":
		contentType, err := parser.GetStringByName("contentType")
		if errors.Is(err, script.ErrMissingKwarg) {
			contentType = "application/json"
		} else if err != nil {
			return nil, err
		}
		data, err := parser.GetStringByName("data")
		if err != nil {
			return nil, err
		}
		resp, err = http.Post(requestURL, contentType, bytes.NewBufferString(data))
		if err != nil {
			return nil, err
		}
	default:
		resp, err = http.Get(requestURL)
		if err != nil {
			return nil, err
		}
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if file != "" {
		err = ioutil.WriteFile(file, respBytes, 0644)
		if err != nil {
			return nil, err
		}
	}
	return string(respBytes), nil
}
