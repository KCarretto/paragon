package file_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/assert"
	"github.com/kcarretto/paragon/pkg/script/stdlib/file"
	"github.com/kcarretto/paragon/pkg/script/stdlib/sys"
	"github.com/stretchr/testify/require"
)

func execScript(t *testing.T, name string, code string) error {
	return script.New(name, bytes.NewBufferString(code),
		assert.Include(),
		file.Include(),
		sys.Include(),
	).
		Exec(context.Background())
}

const codeTestOperations = `
def test_sys_openFile(fileName):
	f, err = sys.openFile(fileName)
	assert.noError(err)
	return f

def test_sys_exec(cmd):
	out, err = sys.exec(cmd)
	assert.noError(err)
	print("============= Exec =============")
	print("Ran Command: "+cmd)
	print(out)
	print("============= ---- =============")

def test_name(f):
	print("============= Name =============")
	print(file.name(f))
	print("============= ---- =============")

def test_write(f, fileContent):
	path = file.name(f)

	err = file.write(f, fileContent)
	assert.noError(err)

	f2 = test_sys_openFile(path)

	content, err = file.content(f2)
	assert.noError(err)

	if content != fileContent:
		fail("Failed to write expected content to file", fileContent, content)

def test_content(f, fileContent):
	print("============= Content =============")
	content, err = file.content(f)
	assert.noError(err)
	print(content)
	print("============= ---- =============")

def test_move(f, dstPath):
	srcPath = file.name(f)
	err = file.move(f, dstPath)
	assert.noError(err)

	oldF = test_sys_openFile(srcPath)
	content, err = file.content(f)
	assert.noError(err)

	oldContent, err = file.content(oldF)
	assert.noError(err)

def test_copy(f1, dstPath):
	srcPath = file.name(f1)
	content1, err = file.content(f1)
	assert.noError(err)

	f2 = test_sys_openFile(dstPath)
	err = file.copy(f1, f2)
	assert.noError(err)

	content2, err = file.content(f2)
	assert.noError(err)

	if content1 != content2:
		fail("Copied file content differed from original content", content1, content2)
	return f2

def test_remove(path):
	f = test_sys_openFile(path)
	err = file.remove(f)
	assert.noError(err)

# TODO: Requires root / to know what the operating user is
#def test_chown(f):
#	pass
#	err = file.chown(f, "root", "root")
#	assert.noError(err)

def test_chmod(f):
	err = file.chmod(f, "777")
	assert.noError(err)

def main():
	os = sys.detectOS()

	prefix = "/tmp/paragon_test/" if os != "windows" else "C:\\Windows\\"
	cmd = "ls -al /" if os != "windows" else "dir"

	fileName = prefix + "path/to/file"
	newPath = prefix + "new_path/to/file"
	newNewPath = prefix + "new_new_path/to/file"
	fileContent = "boop"

	print("Opening File: "+fileName)
	f = test_sys_openFile(fileName)
	test_name(f)
	test_write(f, fileContent)
	test_content(f, fileContent)
	test_move(f, newPath)
	f = test_copy(f, newNewPath)
	test_remove(newPath)
	#test_chown(f)
	test_chmod(f)
	test_sys_exec(cmd)
`

func TestOperations(t *testing.T) {
	err := execScript(t, "operations_test_script", codeTestOperations)
	require.NoError(t, err, "script failed execution")
}
