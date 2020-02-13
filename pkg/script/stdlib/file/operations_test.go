package file_test

import (
	"bytes"
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/file"
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
	return script.New(name, bytes.NewBufferString(code), file.Include(), testLib(f)).
		Exec(context.Background())
}

const codeTestMove = `
def main():
	f = testlib.f()
	file.move(f, "/path/to/dst")
`

func TestMove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dst := "/path/to/dst"

	f := NewMockFile(ctrl)
	f.EXPECT().Move(gomock.Eq(dst)).Return(nil)

	err := execScript(t, f, "move_test", codeTestMove)
	require.NoError(t, err, "script failed execution")
}
