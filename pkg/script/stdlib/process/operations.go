package process

import (
	"github.com/kcarretto/paragon/pkg/script"

	psutil "github.com/shirou/gopsutil/process"
)

// Kill a process (using SIGKILL).
//
//go:generate go run ../gendoc.go -lib process -func kill -param proc@Process -retval err@Error -doc "Kill a process (using SIGKILL)."
//
// @callable:	process.Kill
// @param:		proc	@Process
// @retval:		err 	@Error
//
// @usage:		err = process.Kill(proc)
func Kill(proc Process) error {
	p, err := psutil.NewProcess(proc.Pid)
	if err != nil {
		return nil
	}
	return p.Kill()
}

func kill(parser script.ArgParser) (script.Retval, error) {
	p, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}
	return Kill(p), nil
}

// Name gets the name of the passed process.
//
//go:generate go run ../gendoc.go -lib process -func name -param proc@Process -retval name@String -retval err@Error -doc "Name gets the name of the passed process."
//
// @callable:	process.Name
// @param:		proc	@Process
// @retval:		name 	@String
// @retval:		err 	@Error
//
// @usage:		proc_name = process.Name(proc)
func Name(proc Process) string {
	return proc.Name
}

func name(parser script.ArgParser) (script.Retval, error) {
	p, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}
	return Name(p), nil
}
