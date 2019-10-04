package script_test

import (
	"testing"

	"github.com/kcarretto/paragon/script"
)

func convertTestString(argParse script.ArgParser) (script.Retval, error) {
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
	// newFunc := script.Func(convertTestString)
	// i := script.NewInterpreter()
	// l := script.Library{"my_func": newFunc}
	// i.AddLibrary("mylib", l)

	// script := script.New("myscript", []byte(myconvertscript))
	// logger, _ := zap.NewDevelopment()
	// err := i.Exec(context.Background(), logger, script)
	// if err != nil {
	// 	t.Error("Error executing test: ", err)
	// }
	// correctData := "[myscript] [True, 1, 2, 3, 4.4, 5.5, \"1\", {\"1\": \"1\"}, {\"1\": \"1\"}, None]\n"
	// t.Log(logger)                         // Borked
	// require.Equal(t, output.String(), correctData) // Borked
}
