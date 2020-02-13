package file

import (
	"io"
	"io/ioutil"
	"os/user"
	"strconv"

	"github.com/kcarretto/paragon/pkg/script"
)

// Move a file to the desired location.
//
// @callable:	file.move
// @param:		file	@File
// @param:		dstPath @str
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

	return nil, Move(f, dstPath)
}

// Name returns file's basename.
//
// @callable: 	file.name
// @param: 		file @File
// @retval:		name @str
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
// @callable: 	file.content
// @param: 		file 	@File
// @retval:		content @str
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

	return Content(f)
}

// Write sets the file's content.
//
// @callable: 	file.write
// @param: 		file 	@File
// @param:		content @str
// @retval:		err 	@Error
//
// @usage: 		err = file.write(f, "New\n\tFile\n\t\tContent")
func Write(file Type, content string) error {
	_, err := Type.Write(file, []byte(content))
	return err
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

	return nil, Write(f, content)
}

// Copy the file's content into another file.
//
// @callable: 	file.copy
// @param: 		src 	@File
// @param:		dst		@File
// @retval:		err 	@Error
//
// @usage: 		err = file.copy(f1, f2)
func Copy(src Type, dst Type) error {
	_, err := io.Copy(dst, src)
	return err
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

	return nil, Copy(src, dst)
}

// Remove the file. It will become unuseable after calling this operation.
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

	return nil, Remove(f)
}

// Chown modifies the file's ownership metadata. Passing an empty string for either the username
// or group parameter will result in a no-op. For example, file.chown(f, "", "new_group") will
// change the file's group ownership to "new_group" but will not affect the file's user ownership.
//
// @callable: 	file.chown
// @param: 		file 		@File
// @param:		username	@str
// @param:		group		@str
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

	return Type.Chown(file, int(uid), int(gid))
}

func chown(parser script.ArgParser) (script.Retval, error) {
	f, err := parseFileParam(parser, 0)
	if err != nil {
		return nil, err
	}

	username, _ := parser.GetString(1)
	group, _ := parser.GetString(2)

	return nil, Chown(f, username, group)
}
