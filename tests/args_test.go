package tests

import (
	"errors"
	"io"
	"testing"

	"github.com/kcarretto/paragon/interpreter"
)

func aTestFunc(argParse interpreter.ArgParser, output io.Writer) (interpreter.Retval, error) {
	stringVar, err := argParse.GetString(0)
	if err != nil {
		return nil, err
	}
	intVar, err := argParse.GetInt(1)
	if err != nil {
		return nil, err
	}
	boolVar, err := argParse.GetBool(2)
	if err != nil {
		return nil, err
	}
	if stringVar != "test" {
		return nil, errors.New("String var was wrong")
	}
	if intVar != 1337 {
		return nil, errors.New("Int var was wrong")
	}
	if boolVar != true {
		return nil, errors.New("Bool var was wrong")
	}
	return nil, nil
}

type Output struct{}

func (o *Output) Write(p []byte) (int, error) {
	return len(p), nil
}

const myscript string = `
load("mylib", "my_func")

def main():
	print("test")
`

func TestArgParse(t *testing.T) {
	newFunc := interpreter.Func(aTestFunc)
	i := interpreter.New()
	l := interpreter.Library{"my_func": newFunc}
	i.AddLibrary("mylib", l)

	script := interpreter.NewScript("myscript", []byte(myscript))
	output := &Output{}
	err := i.Execute(script, output)
	if err != nil { // erroring here
		t.Error("Error executing test: ", err)
	}
}
