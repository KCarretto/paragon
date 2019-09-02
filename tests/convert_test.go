package tests

import (
	"fmt"
	"io"
	"testing"

	"github.com/kcarretto/paragon/interpreter"
)

func convertTestString(argParse interpreter.ArgParser, output io.Writer) (interpreter.Retval, error) {
	test := make([]interface{}, 9)
	test[0] = true
	test[1] = int(1)
	test[2] = int64(2)
	test[3] = uint64(3)
	test[4] = float32(4.4)
	test[5] = float64(5.5)
	test[6] = "1"
	test[7] = map[string]string{"1": "1"}
	test[8] = nil
	return test, nil
}

const myconvertscript string = `
load("mylib", "my_func")

def main():
	a = my_func()
	print(a)
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
	fmt.Println(output.String())
	t.Error("wa")
}
