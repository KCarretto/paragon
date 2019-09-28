package sys

import (
	"io/ioutil"

	"github.com/kcarretto/paragon/script"
)

// ReadFile uses ioutil.ReadFile to read an entire file's contents.
//
// @param file: A string for the path of the file.
//
// @return (fileContents, nil) iff success; (nil, err) o/w
func ReadFile(parser script.ArgParser) (script.Retval, error) {
	file, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadFile(file)
	return string(result), err
}

// WriteFile uses ioutil.WriteFile to write an entire file's contents, perms are set to 0644.
//
// @param file: A string for the path of the file.
// @param content: A string for the content of the file to be written to.
//
// @return (nil, nil) iff success; (nil, err) o/w
func WriteFile(parser script.ArgParser) (script.Retval, error) {
	file, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	content, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}
	// if you fucks want different perms use chmod
	err = ioutil.WriteFile(file, []byte(content), 0644)
	return nil, err
}
