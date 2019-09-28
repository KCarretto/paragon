package sys

import (
	"os"

	"github.com/kcarretto/paragon/script"
)

// Remove uses os.Rename to remove a file/folder. WARNING: basically works like rm -rf.
//
// @param file: A string for the path of the file.
//
// @return (nil, nil) iff success; (nil, err) o/w
func Remove(parser script.ArgParser) (script.Retval, error) {
	file, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	err = os.RemoveAll(file)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
