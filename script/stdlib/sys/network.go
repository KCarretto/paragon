package sys

import (
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

// Connections uses net.ConnectionsPid to get all the conenctions opened by a process, given a connection protocol.
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
