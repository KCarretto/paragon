package tests

import (
	"io"
	"testing"

	"github.com/kcarretto/paragon/interpreter"
)

func convertTestString(argParse interpreter.ArgParser, output io.Writer) (interpreter.Retval, error) {
	return "", nil
}

const myconvertscript string = `
load("mylib", "my_func")

def main():
	my_func()
`

func TestConvert(t *testing.T) {
	newFunc := interpreter.Func(convertTestString)
	i := interpreter.New()
	l := interpreter.Library{"my_func": newFunc}
	i.AddLibrary("mylib", l)

	script := interpreter.NewScript("myscript", []byte(myconvertscript))
	output := &Output{}
	err := i.Execute(script, output)
	if err != nil {
		t.Error("Error executing test: ", err)
	}
	t.Error("no test written")
}
