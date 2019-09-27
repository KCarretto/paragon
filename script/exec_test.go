package script_test

import (
	"context"
	"testing"

	"github.com/kcarretto/paragon/script"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestExecOutput(t *testing.T) {
	py := script.NewInterpreter()
	script := script.New("test_execute", []byte(testScriptContent))

	err := py.Exec(context.Background(), zap.NewNop(), script)
	require.NoError(t, err)

	// TODO: Check output of logger

	// result, err := ioutil.ReadAll(reader)
	// require.NoError(t, err)

	// require.Equal(t, testScriptOutput, string(result))
}
