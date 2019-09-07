package interpreter_test

import (
	"testing"

	"github.com/kcarretto/paragon/interpreter"
	"github.com/stretchr/testify/require"
)

const testScriptContent = `
print("loading")

def count():
	nums = [str(x) for x in range(10)]
	print(",".join(nums))

def main():
    count()
`
const testScriptOutput = "[test_execute] loading\n[test_execute] 0,1,2,3,4,5,6,7,8,9\n"

func TestNewScript(t *testing.T) {
	script := interpreter.NewScript("test_execute", []byte(testScriptContent))

	val := make([]byte, len(testScriptContent))
	n, err := script.Read(val)
	require.NoError(t, err)
	require.Equal(t, len(testScriptContent), n)
	require.Equal(t, testScriptContent, string(val))
}