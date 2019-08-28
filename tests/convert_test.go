package tests

import (
	"io"
	"testing"

	"github.com/kcarretto/paragon/interpreter"
)

func convertTestFunc(argParse interpreter.ArgParser, output io.Writer) (interpreter.Retval, error) {
	return nil, nil
}

func TestConvert(t *testing.T) {
	newFunc := interpreter.Func(convertTestFunc)
	i := interpreter.New()
	l := interpreter.Library{"my_func": newFunc}
	i.AddLibrary("mylib", l)

	script := interpreter.NewScript("myscript", []byte(myscript))
	output := &Output{}
	err := i.Execute(script, output)
	if err != nil {
		t.Error("Error executing test: ", err)
	}
	t.Error("no test written")
}
