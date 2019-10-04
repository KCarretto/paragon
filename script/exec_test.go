package script_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/kcarretto/paragon/script"
	"github.com/stretchr/testify/require"
)

func TestExecOutput(t *testing.T) {
	code := script.New("test_execute", bytes.NewBufferString(testScriptContent))

	err := code.Exec(context.Background())
	require.NoError(t, err)

	// TODO: Check output of logger

	// result, err := ioutil.ReadAll(reader)
	// require.NoError(t, err)

	// require.Equal(t, testScriptOutput, string(result))
}
