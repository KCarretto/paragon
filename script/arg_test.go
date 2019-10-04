package script_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/kcarretto/paragon/script"
)

func aTestFunc(argParse script.ArgParser) (script.Retval, error) {
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

const myscript string = `
load("mylib", "my_func")

def main():
	my_func("test", 1337, True)
`

func TestArgParse(t *testing.T) {
	newFunc := script.Func(aTestFunc)
	i := script.NewInterpreter()
	l := script.Library{"my_func": newFunc}
	i.AddLibrary("mylib", l)

	code := script.New("myscript", bytes.NewBufferString(myscript))
	err := code.Exec(context.Background())
	if err != nil {
		t.Error("Error executing test: ", err)
	}
}
