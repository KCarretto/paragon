package interpreter_test

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/kcarretto/paragon/interpreter"
	"github.com/stretchr/testify/require"
)

func TestExecOutput(t *testing.T) {
	py := interpreter.New()
	script := interpreter.NewScript("test_execute", []byte(testScriptContent))

	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()

		err := py.Execute(script, writer)
		require.NoError(t, err)
	}()

	result, err := ioutil.ReadAll(reader)
	require.NoError(t, err)

	require.Equal(t, testScriptOutput, string(result))
}
