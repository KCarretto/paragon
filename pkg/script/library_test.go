package script_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/kcarretto/paragon/pkg/script"
	"github.com/stretchr/testify/require"
)

const testLibraryScriptContent = `
print(testlib.do())
`
const testLibraryScriptOutput = "Success!\n"

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
	lib := script.Library(map[string]script.Func{
		"do": getDo(),
	})

	var output bytes.Buffer
	code := script.New("test_library", bytes.NewBufferString(testLibraryScriptContent), script.WithLibrary("testlib", lib), script.WithOutput(&output))
	err := code.Exec(context.Background())
	require.NoError(t, err)
	require.Equal(t, testLibraryScriptOutput, output.String())
}
