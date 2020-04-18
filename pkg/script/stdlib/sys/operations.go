package sys

import (
	"fmt"
	os_exec "os/exec"
	"strings"

	"github.com/kcarretto/paragon/pkg/script"
	libfile "github.com/kcarretto/paragon/pkg/script/stdlib/file"
	libnet "github.com/kcarretto/paragon/pkg/script/stdlib/net"
	libproc "github.com/kcarretto/paragon/pkg/script/stdlib/process"

	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/afero"
)

//go:generate go run ../gendoc.go -lib sys -func file -param path@String -retval f@File -doc "Prepare a descriptor for a file on the system. The descriptor may be used with the file library."
func file(parser script.ArgParser) (script.Retval, error) {
	path, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}

	return &libfile.File{
		Path: path,
		Fs:   afero.NewOsFs(),
	}, nil
}

//go:generate go run ../gendoc.go -lib sys -func exec -param executable@String -param disown@?Bool -retval output@String -retval err@Error -doc "Exec uses the os/exec.command to execute the passed executable/params. Disown will optionally spawn a process but prevent it's output from being returned."
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

	argv := strings.Fields(executable)
	if len(argv) == 0 {
		return "", fmt.Errorf("exec expected args but got none")
	}

	// TODO: Possibly use bin.Process to return state about the spawned process to the user.
	bin := os_exec.Command(argv[0], argv[1:]...)
	if disown {
		return script.WithError("", bin.Start()), nil
	}

	result, err := bin.CombinedOutput()
	return script.WithError(string(result), err), nil
}

//go:generate go run ../gendoc.go -lib sys -func connections -param parent@?Process -retval connections@[]Connection -doc "Connections uses the gopsutil/net to get all connections created by a process (or all by default)."
func connections(parser script.ArgParser) (script.Retval, error) {
	ppid := int32(0)
	if parent, err := libproc.ParseParam(parser, 0); err == nil {
		ppid = parent.Pid
	}

	connections, err := net.ConnectionsPid("all", ppid)
	var conns []libnet.Connection
	if err != nil {
		// TODO: Log error
		return conns, nil
	}

	for _, conn := range connections {
		conns = append(conns, libnet.Connection{
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

//go:generate go run ../gendoc.go -lib sys -func processes -retval procs@[]Process -doc "Processes uses the gopsutil/process to get all processes."
func processes(parser script.ArgParser) (script.Retval, error) {
	procs := []libproc.Process{}
	pids, err := process.Pids()
	if err != nil {
		// TODO: Log error
		return procs, nil
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

		procs = append(procs, libproc.Process{
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

// Files uses the ioutil.ReadDir to get all files in a given path.
//
// TODO: Remove spaces before go:generate
//  go:generate go run ../gendoc.go -lib sys -func files -retval files@[]File -retval err@Error -doc "Files uses the ioutil.ReadDir to get all files in a given path."
//
// @callable:	sys.files
// @retval:		files 		@[]File
// @retval:		err 		@Error
//
// @usage:		files, err = sys.files()
// func Files(path string) (files []filelib.File, err error) {
// 	fileInfos, err := ioutil.ReadDir(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, info := range fileInfos {
// 		if !info.IsDir() {
// 			f, err := OpenFile(os_path.Join(path, info.Name()))
// 			if err != nil {
// 				return nil, err
// 			}
// 			files = append(files, f)
// 		}
// 	}
// 	return files, nil
// }

// func files(parser script.ArgParser) (script.Retval, error) {
// 	path, err := parser.GetString(0)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return script.WithError(Files(path)), nil
// }
