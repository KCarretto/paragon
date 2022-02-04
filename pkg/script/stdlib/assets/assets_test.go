package assets_test

import (
	"io/ioutil"
	"testing"

	"github.com/kcarretto/paragon/pkg/script/stdlib/assets"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func writeFileForTest(fs afero.Fs, filename string, content string) error {
	f, err := fs.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func TestTarGZBundleConsistent(t *testing.T) {
	fs := afero.NewMemMapFs()
	file1Name, file1Content := "file1", "boop"
	file2Name, file2Content := "file2", "shmoop"

	err := writeFileForTest(fs, file1Name, file1Content)
	require.NoError(t, err, "failed to create test file")

	err = writeFileForTest(fs, file2Name, file2Content)
	require.NoError(t, err, "failed to create test file")

	f1, err := fs.Open(file1Name)
	require.NoError(t, err, "failed to open file")

	f2, err := fs.Open(file2Name)
	require.NoError(t, err, "failed to open file")

	targzbundlr := assets.TarGZBundler{}
	err = targzbundlr.Bundle(
		assets.NamedReader{
			Reader: f1,
			Name:   file1Name,
		},
		assets.NamedReader{
			Reader: f2,
			Name:   file2Name,
		},
	)
	require.NoError(t, err, "failed to bundle files")

	newFS, err := targzbundlr.FileSystem()
	require.NoError(t, err, "failed to open unbundle the files into a filesystem")

	newF1, err := newFS.Open(file1Name)
	require.NoError(t, err, "failed to open untargz'd file")

	newF2, err := newFS.Open(file2Name)
	require.NoError(t, err, "failed to open untargz'd file")

	newf1Content, err := ioutil.ReadAll(newF1)
	require.NoError(t, err, "failed to read file")

	newf2Content, err := ioutil.ReadAll(newF2)
	require.NoError(t, err, "failed to read file")

	require.Equal(t, file1Content, string(newf1Content), "invalid file content")
	require.Equal(t, file2Content, string(newf2Content), "invalid file content")
}
