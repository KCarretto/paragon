package sys

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/kcarretto/paragon/script"
)

// Exec uses os.Command to execute the passed string
//
// @param cmd: A string for execution.
//
// @return (stdoutStderr, nil) iff success; (nil, err) o/w
func Exec(parser script.ArgParser) (script.Retval, error) {
	cmd, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	argv := strings.Fields(cmd)
	if len(argv) == 0 {
		return nil, errors.New("exec expected args but got none")
	}
	bin := exec.Command(argv[0], argv[1:]...)
	result, err := bin.CombinedOutput()
	return string(result), err
}
