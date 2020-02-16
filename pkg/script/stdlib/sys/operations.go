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
//
// @callable: 	sys.openFile
// @param: 		path 	@string
// @retval:		file 	@File
// @retval:		err 	@Error
//
// @usage: 		f, err = sys.openFile("/usr/bin/malware")
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

	retVal, retErr := OpenFile(path)
	return script.WithError(retVal, retErr), nil
}

// DetectOS uses the GOOS variable to determine the OS.
//
// @callable:	sys.detectOS
// @retval:		os  	@string
//
// @usage:		os = sys.detectOS()
func DetectOS() string {
	switch runtime.GOOS {
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	case "darwin":
		return "darwin"
	case "freebsd", "netbsd", "openbsd", "dragonfly":
		return "bsd"
	case "android":
		return "android"
	case "solaris":
		return "solaris"
	default:
		return "other"
	}
}

func detectOS(parser script.ArgParser) (script.Retval, error) {
	return DetectOS(), nil
}

// Exec uses the os/exec.command to execute the passed executable/params.
//
// @callable:	sys.exec
// @param:		executable	@string
// @param:		disown		@?bool
// @retval:		output 		@string
// @retval:		err 		@Error
//
// @usage:		output = sys.exec("/usr/sbin/nginx", disown=True)
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

	retVal, retErr := Exec(executable, disown)
	return script.WithError(retVal, retErr), nil
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
