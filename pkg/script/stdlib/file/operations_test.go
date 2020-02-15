package file_test

import (
	"bytes"
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/file"
	"github.com/kcarretto/paragon/pkg/script/stdlib/sys"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -destination=operations.gen_test.go -package=file_test github.com/kcarretto/paragon/pkg/script/stdlib/file File

func testLib(f file.File) script.Option {
	return script.WithLibrary("testlib", script.Library(map[string]script.Func{
		"f": script.Func(func(parser script.ArgParser) (script.Retval, error) {
			return file.New(f), nil
		})}),
	)
}

func execScript(t *testing.T, f file.File, name string, code string) error {
	return script.New(name, bytes.NewBufferString(code), testLib(f),
		file.Include(),
		sys.Include(),
	).
		Exec(context.Background())
}

const codeTestMove = `
def main():
	f = sys.openFile("path/to/src")
	f2 = testlib.f()
	file.move(f, "dst")
	file.move(f2, "dst")

`

func TestMove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dst := "dst"

	f := NewMockFile(ctrl)
	f.EXPECT().Move(gomock.Eq(dst)).Return(nil)

	err := execScript(t, f, "move_test", codeTestMove)
	require.NoError(t, err, "script failed execution")
}
