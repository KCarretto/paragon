package sys

import (
	"strconv"

	"github.com/kcarretto/paragon/script"
	"github.com/shirou/gopsutil/process"
)

// Processes uses gopsutil.process.Pids to get all pids for a box and then makes them into Process map structs.
//
// BUG(cictrone): It seems this functions is much much much slower on Darwin.
//
// The Process map is defined as:
//  map[string]string{
//   "pid":     PID of the process,
//   "ppid":    PID of the parent process,
//   "name":    Name of the process,
//   "user":    User who ran/owns the process,
//   "status":  Status of the process,
//   "cmdLine": Command line argumenest for the process,
//   "exe":     Name of executable that started the process,
//   "tty":     The tty/pty of the process ,
//  }
//
// @return (processes, nil) iff success; (nil, err) o/w
func Processes(parser script.ArgParser) (script.Retval, error) {
	pids, err := process.Pids()
	if err != nil {
		return nil, err
	}
	var procs []map[string]string
	for _, pid := range pids {
		proc, err := process.NewProcess(pid)
		if err != nil {
			// if a process dies between getting the list and now, its okay
			continue
		}
		procPidStr := strconv.FormatInt(int64(proc.Pid), 10)
		procPpid, _ := proc.Ppid()
		procPpidStr := strconv.FormatInt(int64(procPpid), 10)
		procName, _ := proc.Name()
		procUser, _ := proc.Username()
		procStatus, _ := proc.Status()
		procCmdLine, _ := proc.Cmdline()
		procExe, _ := proc.Exe()
		procTerminal, _ := proc.Terminal()

		procMap := map[string]string{
			"pid":     procPidStr,
			"ppid":    procPpidStr,
			"name":    procName,
			"user":    procUser,
			"status":  procStatus,
			"cmdLine": procCmdLine,
			"exe":     procExe,
			"tty":     procTerminal,
		}

		procs = append(procs, procMap)
	}
	return procs, nil
}

// Kill uses gopsutil.process.Kill to kill and passed process pid.
//
// @param pid: A string of the process pid to be killed.
//
// @return (nil, nil) iff success; (nil, err) o/w
func Kill(parser script.ArgParser) (script.Retval, error) {
	pid, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	pid64, err := strconv.ParseInt(pid, 10, 32)
	if err != nil {
		return nil, err
	}
	proc, err := process.NewProcess(int32(pid64))
	if err != nil {
		return nil, err
	}
	return nil, proc.Kill()
}
