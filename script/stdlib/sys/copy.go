package sys

import (
	"io/ioutil"

	"github.com/kcarretto/paragon/script"
)

// Copy uses ioutil.ReadFile and ioutil.WriteFile to copy a file from source to destination.
//
// @param srcFile: A string for the path of the source file.
// @param dstFile: A string for the path of the destination file.
//
// @return (nil, nil) iff success; (nil, err) o/w
func Copy(parser script.ArgParser) (script.Retval, error) {
	srcFile, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	dstFile, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}

	result, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(dstFile, result, 0777)

	return nil, err
}
