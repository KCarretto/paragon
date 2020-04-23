package file

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/kcarretto/paragon/pkg/script"
	"golang.org/x/crypto/sha3"

	"github.com/spf13/afero"
)

//go:generate go run ../gendoc.go -lib file -func move -param f@File -param dstPath@String -retval err@Error -doc "Move a file to the desired location."
func move(parser script.ArgParser) (script.Retval, error) {
	fd, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}

	dstPath, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}

	if err := rename(fd, fd.Path, dstPath); err != nil {
		return err, nil
	}

	fd.Path = dstPath
	return nil, nil
}

//go:generate go run ../gendoc.go -lib file -func hash -param f@File -retval digest@String -doc "The SHA3-256 hash of the passed file (base64 encoded)."
func hash(parser script.ArgParser) (script.Retval, error) {
	fd, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}
	src, err := fd.OpenFile(fd.Path, os.O_RDONLY, 0755)
	if err != nil {
		return script.WithError(nil, fmt.Errorf("failed to open file: %w", err)), nil
	}
	defer src.Close()
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return script.WithError(nil, fmt.Errorf("failed to read file: %w", err)), nil
	}
	digest := sha3.Sum256(data)
	return base64.StdEncoding.EncodeToString(digest[:]), nil
}

//go:generate go run ../gendoc.go -lib file -func exists -param f@File -retval exists@Bool -doc "The boolean value on if the file exists or not."
func exists(parser script.ArgParser) (script.Retval, error) {
	fd, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(fd.Path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

//go:generate go run ../gendoc.go -lib file -func name -param f@File -retval name@String -doc "The name or path used to open the file."
func name(parser script.ArgParser) (script.Retval, error) {
	fd, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}
	return fd.Path, nil
}

//go:generate go run ../gendoc.go -lib file -func content -param f@File -retval content@String -doc "Read and return the file's contents."
func content(parser script.ArgParser) (script.Retval, error) {
	fd, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}

	f, err := fd.Open(fd.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	return string(content), err
}

//go:generate go run ../gendoc.go -lib file -func write -param f@File -param content@String -doc "Write sets the file's content, overwriting any previous value. It creates the file if it does not yet exist."
func write(parser script.ArgParser) (script.Retval, error) {
	fd, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}

	content, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}

	f, err := upsertFile(fd)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return nil, err
}

//go:generate go run ../gendoc.go -lib file -func copy -param src@File -param dst@File -retval err@Error -doc "Copy the file's content into a destination file, overwriting any previous value. It creates the destination file if it does not yet exist."
func copy(parser script.ArgParser) (script.Retval, error) {
	srcFD, err := ParseParam(parser, 0)
	if err != nil {
		return nil, fmt.Errorf("src file: %w", err)
	}

	dstFD, err := ParseParam(parser, 1)
	if err != nil {
		return nil, fmt.Errorf("dst file: %w", err)
	}

	src, err := srcFD.OpenFile(srcFD.Path, os.O_RDONLY, 0755)
	if err != nil {
		return script.WithError(nil, fmt.Errorf("failed to open src file: %w", err)), nil
	}
	defer src.Close()

	dst, err := upsertFile(dstFD)
	if err != nil {
		return script.WithError(nil, fmt.Errorf("failed to open dst file: %w", err)), nil
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return script.WithError(nil, fmt.Errorf("failed to copy files: %w", err)), nil
	}

	return nil, nil
}

//go:generate go run ../gendoc.go -lib file -func remove -param f@File -retval err@Error -doc "Remove the file"
func remove(parser script.ArgParser) (script.Retval, error) {
	fd, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}

	return fd.RemoveAll(fd.Path), nil
}

//go:generate go run ../gendoc.go -lib file -func chmod -param f@File -param mode@String -doc "Chmod modifies the file's permission metadata. The strong passed is expected to be an octal representation of what os.FileMode you wish to set file to have (i.e. '0755')."
func chmod(parser script.ArgParser) (script.Retval, error) {
	fd, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}

	modeStr, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}

	mode, err := strconv.ParseInt(modeStr, 8, 32)
	if err != nil {
		return nil, err
	}

	return nil, fd.Chmod(fd.Path, os.FileMode(mode))
}

//go:generate go run ../gendoc.go -lib file -func drop -param src@File -param dst@File -param perms@?String -retval err@Error -doc "Drop will:\n\t1. Copy a given file to a tempfile on disk\n\t2. Optionally set the permissions The default perms are '0755'.\n\t3. Move it to a given destination\n\t4. Clean up the temp file created."
func drop(parser script.ArgParser) (script.Retval, error) {
	srcFD, err := ParseParam(parser, 0)
	if err != nil {
		return nil, fmt.Errorf("src file: %w", err)
	}

	dstFD, err := ParseParam(parser, 1)
	if err != nil {
		return nil, fmt.Errorf("dst file: %w", err)
	}

	var perms int64 = 0755
	if permStr, err := parser.GetString(2); err == nil {
		perms, err = strconv.ParseInt(permStr, 8, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid perms: %w", err)
		}
	}

	// Open Source File
	src, err := srcFD.Open(srcFD.Path)
	if err != nil {
		return fmt.Errorf("failed to open src file: %w", err), nil
	}

	// Create TempFile on Destination Filesystem
	dstFS := afero.Afero{
		Fs: dstFD.Fs,
	}
	tmp, err := dstFS.TempFile(dstFS.GetTempDir(""), strconv.Itoa(rand.Int()))
	if err != nil {
		return fmt.Errorf("failed to open temp file on destination filesystem: %w", err), nil
	}
	tmpPath := tmp.Name()
	defer func() {
		tmp.Close()
		dstFS.Remove(tmpPath)
		dstFS.RemoveAll(tmpPath)
	}()

	// Copy Src to Temp File
	if _, err := io.Copy(tmp, src); err != nil {
		return fmt.Errorf("failed to copy src into temp file: %w", err), nil
	}

	// Move Temp File to Dst
	if err := rename(dstFS, tmpPath, dstFD.Path); err != nil {
		return fmt.Errorf("failed to replace dst file with temp file: %w", err), nil
	}

	// Update DST File permissions
	if err := dstFS.Chmod(dstFD.Path, os.FileMode(perms)); err != nil {
		return fmt.Errorf("failed to update dst file perms: %w", err), nil
	}

	return nil, nil
}

func rename(fs afero.Fs, src, dst string) error {
	// SSH implementation doesn't rename if dst already exists, so delete it
	if exists, _ := afero.Exists(fs, dst); exists {
		fs.RemoveAll(dst)
		fs.Remove(dst)
	}

	dir := filepath.Dir(dst)
	if exists, err := afero.DirExists(fs, dir); !exists || err != nil {
		fs.MkdirAll(dir, 0755)
	}

	if err := fs.Rename(src, dst); err != nil {
		return err
	}

	return nil
}

func upsertFile(fd *File) (afero.File, error) {
	dir := filepath.Dir(fd.Path)
	if exists, err := afero.DirExists(fd, dir); !exists || err != nil {
		fd.MkdirAll(dir, 0755)
	}

	f, err := fd.OpenFile(fd.Path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	// SSH implementation doesn't pass the flags and returns nil, nil when file does not exist
	if f == nil {
		f, err = fd.Create(fd.Path)
		if err != nil || f == nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
	}

	return f, nil
}
