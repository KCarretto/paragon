package interpreter_test

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kcarretto/paragon/interpreter"
)

const testLibraryScriptContent = `
load("testlib", "do")
print(do())
`
const testLibraryScriptOutput = "[test_library] Success!\n"

func getDo() interpreter.Func {
	return interpreter.Func(func(args interpreter.ArgParser, output io.Writer) (interpreter.Retval, error) {
		// TODO: Use logger to write to provided stream
		// logger, err := zap.NewDevelopment()
		// if err != nil {
		// 	logger := zap.NewExample()
		// }

		return "Success!", nil
	})
}
func TestLibraryMethod(t *testing.T) {
	py := interpreter.New()
	script := interpreter.NewScript("test_library", []byte(testLibraryScriptContent))

	py.AddLibrary("testlib", interpreter.Library{
		"do": getDo(),
	})

	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()

		err := py.Execute(script, writer)
		require.NoError(t, err)
	}()

	result, err := ioutil.ReadAll(reader)
	require.NoError(t, err)

	require.Equal(t, testLibraryScriptOutput, string(result))
}
