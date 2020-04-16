package assets_test

import (
	"io/ioutil"
	"testing"

	"github.com/kcarretto/paragon/pkg/script/stdlib/assets"
	"github.com/spf13/afero"
)

func writeFileForTest(fs *afero.HttpFs, filename string, content string) error {
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
	fs := afero.NewHttpFs(afero.NewMemMapFs())
	file1Name, file1Content := "file1", "boop"
	file2Name, file2Content := "file2", "shmoop"
	if err := writeFileForTest(fs, file1Name, file1Content); err != nil {
		t.Errorf("failed to create test file: %w", err)
	}
	if err := writeFileForTest(fs, file2Name, file2Content); err != nil {
		t.Errorf("failed to create test file: %w", err)
	}
	f1, err := fs.Open(file1Name)
	if err != nil {
		t.Errorf("failed to open file: %w", err)
	}
	f2, err := fs.Open(file2Name)
	if err != nil {
		t.Errorf("failed to open file: %w", err)
	}
	targzbundlr := assets.TarGZBundler{}
	err = targzbundlr.Bundle(f1, f2)
	if err != nil {
		t.Errorf("failed to bundle files: %w", err)
	}
	newFS, err := targzbundlr.FileSystem()
	if err != nil {
		t.Errorf("failed to open unbundle the files into a filesystem: %w", err)
	}
	newF1, err := newFS.Open(file1Name)
	if err != nil {
		t.Errorf("failed to open untargz'd file: %w", err)
	}
	newF2, err := newFS.Open(file2Name)
	if err != nil {
		t.Errorf("failed to open untargz'd file: %w", err)
	}
	newf1Content, err := ioutil.ReadAll(newF1)
	if err != nil {
		t.Errorf("failed to read file: %w", err)
	}
	newf2Content, err := ioutil.ReadAll(newF2)
	if err != nil {
		t.Errorf("failed to read file: %w", err)
	}

	if file1Content != string(newf1Content) {
		t.Errorf("'%s' does not equal '%s'", file1Content, string(newf1Content))
	}
	if file2Content != string(newf2Content) {
		t.Errorf("'%s' does not equal '%s'", file2Content, string(newf2Content))
	}
}
