package sys

import (
	"errors"
	"fmt"
	"os"
	os_exec "os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/file"
	netlib "github.com/kcarretto/paragon/pkg/script/stdlib/net"
	proclib "github.com/kcarretto/paragon/pkg/script/stdlib/process"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

// OpenFile uses os.Open to Open a file.
func OpenFile(filePath string) (file.Type, error) {
	dir := path.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return file.New(nil), fmt.Errorf("failed to create parent directory %q: %w", dir, err)
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return file.New(nil), err
	}
	return file.New(File{
		f,
	}), nil
}

func openFile(parser script.ArgParser) (script.Retval, error) {
	path, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}

	return OpenFile(path)
}

// DetectOS uses the GOOS variable to determine the OS.
//
// @callable:	sys.DetectOS
// @retval:		os  	@String
// @retval:		err 	@Error
//
// @usage:		os = sys.DetectOS()
func DetectOS() (string, error) {
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

func detectOS(parser script.ArgParser) (script.Retval, error) {
	return DetectOS()
}

// Exec uses the os/exec.command to execute the passed executable/params.
//
// @callable:	sys.Exec
// @param:		executable	@String
// @param:		disown		@?Bool
// @retval:		output 		@String
// @retval:		err 		@Error
//
// @usage:		output = sys.Exec("/usr/sbin/nginx", disown=True)
func Exec(executable string, disown bool) (string, error) {
	argv := strings.Fields(executable)
	if len(argv) == 0 {
		return "", errors.New("exec expected args but got none")
	}
	bin := os_exec.Command(argv[0], argv[1:]...)
	if disown {
		bin.Start()
		return "", nil
	}
	result, err := bin.CombinedOutput()
	return string(result), err
}

func exec(parser script.ArgParser) (script.Retval, error) {
	err := parser.RestrictKwargs("disown")
	if err != nil {
		return nil, err
	}
	executable, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	disown, _ := parser.GetBoolByName("disown")
	return Exec(executable, disown)
}

// Connections uses the gopsutil/net to get all connections created by a process (or all by default).
//
// @callable:	sys.Connections
// @param:		parent		@Process
// @retval:		conns 		@[]Connection
// @retval:		err 		@Error
//
// @usage:		connections = sys.Connections()
func Connections(parent *proclib.Process) ([]netlib.Connection, error) {
	ppid := int32(0)
	if parent != nil {
		ppid = parent.Pid
	}
	connections, err := net.ConnectionsPid("all", ppid)
	var conns []netlib.Connection
	if err != nil {
		return conns, err
	}
	for _, conn := range connections {
		conns = append(conns, netlib.Connection{
			Fd:         conn.Fd,
			ConnFamily: conn.Family,
			ConnType:   conn.Type,
			Pid:        conn.Pid,
			LocalIP:    conn.Laddr.IP,
			LocalPort:  conn.Laddr.Port,
			RemoteIP:   conn.Raddr.IP,
			RemotePort: conn.Raddr.Port,
			Status:     conn.Status,
		})
	}
	return conns, nil
}

func connections(parser script.ArgParser) (script.Retval, error) {
	parent, err := proclib.ParseParam(parser, 0)
	if err != nil {
		return Connections(nil)
	}
	return Connections(&parent)
}

// Processes uses the gopsutil/process to get all processes.
//
// @callable:	sys.Processes
// @retval:		procs 		@[]Process
// @retval:		err 		@Error
//
// @usage:		processes = sys.Processes()
func Processes() ([]proclib.Process, error) {
	pids, err := process.Pids()
	var procs []proclib.Process
	if err != nil {
		return procs, err
	}
	for _, pid := range pids {
		proc, err := process.NewProcess(pid)
		if err != nil {
			// if a process dies between getting the list and now, its okay
			continue
		}
		procPpid, _ := proc.Ppid()
		procName, _ := proc.Name()
		procUser, _ := proc.Username()
		procStatus, _ := proc.Status()
		procCmdLine, _ := proc.Cmdline()
		procExe, _ := proc.Exe()
		procTerminal, _ := proc.Terminal()

		procs = append(procs, proclib.Process{
			Pid:     pid,
			PPid:    procPpid,
			Name:    procName,
			User:    procUser,
			Status:  procStatus,
			CmdLine: procCmdLine,
			Exe:     procExe,
			TTY:     procTerminal,
		})
	}
	return procs, nil
}

func processes(parser script.ArgParser) (script.Retval, error) {
	return Processes()
}
