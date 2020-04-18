package file_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/kcarretto/paragon/pkg/script"
	libassert "github.com/kcarretto/paragon/pkg/script/stdlib/assert"
	libenv "github.com/kcarretto/paragon/pkg/script/stdlib/env"
	libfile "github.com/kcarretto/paragon/pkg/script/stdlib/file"
	libsys "github.com/kcarretto/paragon/pkg/script/stdlib/sys"

	"github.com/stretchr/testify/require"
)

func execScript(t *testing.T, name string, code string) error {
	env := &libenv.Environment{}

	return script.New(name, bytes.NewBufferString(code),
		env.Include(),
		libassert.Include(),
		libfile.Include(),
		libsys.Include(),
	).
		Exec(context.Background())
}

const codeTestOperations = `
def test_sys_file(fileName):
	f = sys.file(fileName)
	assert.equal(fileName, file.name(f))

	return f

def test_sys_exec(cmd):
	out, err = sys.exec(cmd)
	assert.noError(err)
	print("============= Exec =============")
	print("Ran Command: "+cmd)
	print(out)
	print("============= ---- =============")

def test_write(f, fileContent):
	file.write(f, fileContent)

def test_move(f, dstPath):
	file.move(f, dstPath)
	assert.equal(dstPath, file.name(f))

def test_content(f, expected):
	assert.equal(expected, file.content(f))

def test_copy(f1, dstPath):
	f2 = test_sys_file(dstPath)
	err = file.copy(f1, f2)
	assert.noError(err)

	return f2

def test_remove(f):
	err = file.remove(f)
	assert.noError(err)

def test_chmod(f):
	file.chmod(f, "777")

def main():
	prefix = "/tmp/paragon_test/" if not env.isWindows() else "C:\\Windows\\"
	cmd = "ls -al /" if not env.isWindows() else "dir"

	fileName = prefix + "path/to/file"
	newPath = prefix + "new_path/to/file"
	newNewPath = prefix + "new_new_path/to/file"
	fileContent = "boop"

	test_sys_exec(cmd)

	f1 = test_sys_file(fileName)

	test_write(f1, fileContent)
	test_content(f1, fileContent)
	test_chmod(f1)

	test_move(f1, newPath)

	f2 = test_copy(f1, newNewPath)
	test_content(f2, fileContent)

	test_remove(f1)
	#test_content(f1, fileContent)
	#test_chown(f1)

`

func TestOperations(t *testing.T) {
	err := execScript(t, "operations_test_script", codeTestOperations)
	require.NoError(t, err, "script failed execution")
}
