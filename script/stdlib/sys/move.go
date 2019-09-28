package sys

import (
	"os"

	"github.com/kcarretto/paragon/script"
)

// Move uses os.Rename to move a file from source to destination.
//
// @param srcFile: A string for the path of the source file.
// @param dstFile: A string for the path of the destination file.
//
// @return (nil, nil) iff success; (nil, err) o/w
func Move(parser script.ArgParser) (script.Retval, error) {
	srcFile, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	dstFile, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}
	err = os.Rename(srcFile, dstFile)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
