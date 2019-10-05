package script_test

import (
	"testing"

	"github.com/kcarretto/paragon/script"
)

const testLibraryScriptContent = `
load("testlib", "do")
print(do())
`
const testLibraryScriptOutput = "[test_library] Success!\n"

func getDo() script.Func {
	return script.Func(func(args script.ArgParser) (script.Retval, error) {
		// TODO: Use logger to write to provided stream
		// logger, err := zap.NewDevelopment()
		// if err != nil {
		// 	logger := zap.NewExample()
		// }

		return script.Retval("Success!"), nil
	})
}
func TestLibraryMethod(t *testing.T) {
	// py := script.NewInterpreter()
	// code := script.New("test_library", []byte(testLibraryScriptContent))

	// py.AddLibrary("testlib", script.Library{
	// 	"do": getDo(),
	// })

	// // reader, writer := io.Pipe()
	// // go func() {
	// // 	defer writer.Close()

	// err := py.Exec(context.Background(), zap.NewNop(), code)
	// require.NoError(t, err)
	// // }()

	// // result, err := ioutil.ReadAll(reader)
	// // require.NoError(t, err)

	// // require.Equal(t, testLibraryScriptOutput, string(result))
}
