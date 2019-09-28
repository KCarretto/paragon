package tests

import (
	"bytes"
	"io"
	"testing"

	"github.com/kcarretto/paragon/interpreter"
	"github.com/stretchr/testify/require"
)

func convertTestString(argParse interpreter.ArgParser, output io.Writer) (interpreter.Retval, error) {
	test := make([]interface{}, 10)
	test[0] = true
	test[1] = int(1)
	test[2] = int64(2)
	test[3] = uint64(3)
	test[4] = float32(4.4)
	test[5] = float64(5.5)
	test[6] = "1"
	test[7] = map[interface{}]interface{}{"1": "1"}
	test[8] = map[string]string{"1": "1"}
	test[9] = nil
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
	correctData := "[myscript] [True, 1, 2, 3, 4.4, 5.5, \"1\", {\"1\": \"1\"}, {\"1\": \"1\"}, None]\n"
	output := bytes.NewBuffer(make([]byte, 0, len(correctData)))
	err := i.Execute(script, output)
	if err != nil {
		t.Error("Error executing test: ", err)
	}
	t.Log(output.String())
	require.Equal(t, output.String(), correctData)
}
