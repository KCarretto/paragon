package sys_test

import (
	"context"
	"testing"
	"bytes"

	"github.com/kcarretto/paragon/script"
	"github.com/kcarretto/paragon/script/stdlib/sys"

	"go.uber.org/zap"
)

const testContent = `
load("sys", "processes")

def main():
	plist = processes()
	print(plist[0])
`

func TestProcesses(t *testing.T) {
	testScript := script.New("test_script", bytes.NewBufferString(testContent), script.WithLibrary("sys", sys.Lib))
	err := testScript.Exec(context.Background())
	if err != nil {
		panic(err)
	}
}
