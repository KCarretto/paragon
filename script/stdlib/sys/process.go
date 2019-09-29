package sys

import (
	"fmt"
	"strconv"

	"github.com/kcarretto/paragon/script"
	"github.com/shirou/gopsutil/process"
)

// Processes uses gopsutil.process.Pids to get all pids for a box and then makes them into Process map structs.
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
//   "terminal":     The tty/pty of the process ,
//  }
//
// @return (processes, nil) iff success; (nil, err) o/w
func Processes(parser script.ArgParser) (script.Retval, error) {
	fmt.Println("entered")
	pids, err := process.Pids()
	if err != nil {
		return nil, err
	}
	fmt.Println("got pids")
	var procs []map[string]string
	fmt.Println(len(pids))
	for i, pid := range pids {
		fmt.Printf("looking through pids: %d\n", i)
		proc, err := process.NewProcess(pid)
		if err != nil {
			fmt.Println("Error!!!")
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

		fmt.Println(procMap)
		procs = append(procs, procMap)
	}
	fmt.Println(procs[0])
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
