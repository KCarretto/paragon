package sys_test

import (
	"context"
	"testing"

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
	intp := script.NewInterpreter()
	intp.AddLibrary("sys", sys.Lib)

	testScript := script.New("test_script", []byte(testContent))
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	if err := intp.Exec(context.Background(), logger, testScript); err != nil {
		panic(err)
	}
}
