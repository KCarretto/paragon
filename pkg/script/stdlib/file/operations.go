package file

import (
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"

	"github.com/kcarretto/paragon/pkg/script"
)

// Move a file to the desired location.
//
//go:generate go run ../gendoc.go -lib file -func move -param file@File -param dstPath@String -retval err@Error -doc "Move a file to the desired location."
//
// @callable:	file.move
// @param:		file	@File
// @param:		dstPath @String
// @retval:		err 	@Error
//
// @usage:		err = file.move(f, "/path/to/dst")
func Move(file Type, dstPath string) error {
	return Type.Move(file, dstPath)
}

func move(parser script.ArgParser) (script.Retval, error) {
	f, err := parseFileParam(parser, 0)
	if err != nil {
		return nil, err
	}

	dstPath, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}

	return Move(f, dstPath), nil
}

// Name returns file's basename.
//
//go:generate go run ../gendoc.go -lib file -func name -param file@File -retval name@String -doc "Name returns file's basename."
//
// @callable: 	file.name
// @param: 		file @File
// @retval:		name @String
//
// @usage:		name = file.name(f)
func Name(file Type) string {
	return Type.Name(file)
}

func name(parser script.ArgParser) (script.Retval, error) {
	f, err := parseFileParam(parser, 0)
	if err != nil {
		return nil, err
	}
	return Name(f), nil
}

// Content returns the file's content.
//
//go:generate go run ../gendoc.go -lib file -func content -param file@File -retval content@String -retval err@Error -doc "Content returns the file's content."
//
// @callable: 	file.content
// @param: 		file 	@File
// @retval:		content @String
// @retval:		err 	@Error
//
// @usage:		content = file.content(f)
func Content(file Type) (string, error) {
	content, err := ioutil.ReadAll(file)
	return string(content), err
}

func content(parser script.ArgParser) (script.Retval, error) {
	f, err := parseFileParam(parser, 0)
	if err != nil {
		return nil, err
	}

	retVal, retErr := Content(f)
	return script.WithError(retVal, retErr), nil
}

// Write sets the file's content.
//
//go:generate go run ../gendoc.go -lib file -func write -param file@File -param content@String -retval err@Error -doc "Write sets the file's content."
//
// @callable: 	file.write
// @param: 		file 	@File
// @param:		content @String
// @retval:		err 	@Error
//
// @usage: 		err = file.write(f, "New\n\tFile\n\t\tContent")
func Write(file Type, content string) error {
	_, err := Type.Write(file, []byte(content))
	if err != nil {
		return err
	}
	return file.Sync()
}

func write(parser script.ArgParser) (script.Retval, error) {
	f, err := parseFileParam(parser, 0)
	if err != nil {
		return nil, err
	}
	content, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}

	return Write(f, content), nil
}

// Copy the file's content into another file.
//
//go:generate go run ../gendoc.go -lib file -func copy -param src@File -param dst@File -retval err@Error -doc "Copy the file's content into another file."
//
// @callable: 	file.copy
// @param: 		src 	@File
// @param:		dst		@File
// @retval:		err 	@Error
//
// @usage: 		err = file.copy(f1, f2)
func Copy(src Type, dst Type) error {
	_, err := io.Copy(dst, src)
	if err != nil {
		return err
	}
	return dst.Sync()
}

func copy(parser script.ArgParser) (script.Retval, error) {
	src, err := parseFileParam(parser, 0)
	if err != nil {
		return nil, err
	}
	dst, err := parseFileParam(parser, 1)
	if err != nil {
		return nil, err
	}

	return Copy(src, dst), nil
}

// Remove the file. It will become unuseable after calling this operation.
//
//go:generate go run ../gendoc.go -lib file -func remove -param file@File -retval err@Error -doc "Remove the file. It will become unuseable after calling this operation."
//
// @callable: 	file.remove
// @param: 		file 	@File
// @retval:		err 	@Error
//
// @usage: 		err = file.remove(f)
func Remove(file Type) error {
	return Type.Remove(file)
}

func remove(parser script.ArgParser) (script.Retval, error) {
	f, err := parseFileParam(parser, 0)
	if err != nil {
		return nil, err
	}

	return Remove(f), nil
}

// Chown modifies the file's ownership metadata. Passing an empty string for either the username
// or group parameter will result in a no-op. For example, file.chown(f, '', 'new_group') will
// change the file's group ownership to 'new_group' but will not affect the file's user ownership.
//
//go:generate go run ../gendoc.go -lib file -func chown -param file@File -param username@String -param group@String -retval err@Error -doc "Chown modifies the file's ownership metadata. Passing an empty string for either the username or group parameter will result in a no-op. For example, file.chown(f, '', 'new_group') will change the file's group ownership to 'new_group' but will not affect the file's user ownership."
//
// @callable: 	file.chown
// @param: 		file 		@File
// @param:		username	@String
// @param:		group		@String
// @retval:		err 		@Error
//
// @usage: 		err = file.chown(f, "root", "sudoers")
func Chown(file Type, username, group string) error {
	var uid int64 = -1
	var gid int64 = -1

	if username != "" {
		userData, err := user.Lookup(username)
		if err != nil {
			return err
		}

		uid, err = strconv.ParseInt(userData.Uid, 10, 32)
		if err != nil {
			return err
		}
	}

	if group != "" {
		groupData, err := user.LookupGroup(group)
		if err != nil {
			return err
		}

		gid, err = strconv.ParseInt(groupData.Gid, 10, 32)
		if err != nil {
			return err
		}
	}

	// No-op if no user or group was specified
	if uid == -1 && gid == -1 {
		return nil
	}

	err := Type.Chown(file, int(uid), int(gid))
	if err != nil {
		return err
	}
	return file.Sync()
}

func chown(parser script.ArgParser) (script.Retval, error) {
	f, err := parseFileParam(parser, 0)
	if err != nil {
		return nil, err
	}

	username, _ := parser.GetString(1)
	group, _ := parser.GetString(2)

	return Chown(f, username, group), nil
}

// Chmod modifies the file's permission metadata. The strong passed is expected to be an octal representation
// of what os.FileMode you wish to set file to have.
//
//go:generate go run ../gendoc.go -lib file -func chmod -param file@File -param mode@String -retval err@Error -doc "Chmod modifies the file's permission metadata. The strong passed is expected to be an octal representation of what os.FileMode you wish to set file to have."
//
// @callable: 	file.chmod
// @param: 		file 		@File
// @param:		mode		@String
// @retval:		err 		@Error
//
// @usage: 		err = file.chmod(f, "0777")
func Chmod(file Type, mode string) error {
	modeInt, err := strconv.ParseInt(mode, 8, 32)
	if err != nil {
		return err
	}
	err = Type.Chmod(file, os.FileMode(modeInt))
	if err != nil {
		return err
	}
	return file.Sync()
}

func chmod(parser script.ArgParser) (script.Retval, error) {
	f, err := parseFileParam(parser, 0)
	if err != nil {
		return nil, err
	}
	modeStr, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}

	return Chmod(f, modeStr), nil
}
