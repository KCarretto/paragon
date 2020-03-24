package assets

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/kcarretto/paragon/pkg/script"
	"github.com/kcarretto/paragon/pkg/script/stdlib/file"
	"github.com/kcarretto/paragon/pkg/script/stdlib/sys"
)

// OpenFile that was packed into the compiled binary. The resulting file does not support many
// operations such as Chown, Write, etc. but you may read it's contents or copy it to another file
// i.e. one opened by ssh or sys.
//
//go:generate go run ../gendoc.go -lib assets -func openFile -param path@String -retval f@File -retval err@Error -doc "OpenFile that was packed into the compiled binary. The resulting file does not support many operations such as Chown, Write, etc. but you may read it's contents or copy it to another file i.e. one opened by ssh or sys."
//
// @callable:	assets.openFile
// @param:		path 	@String
// @retval:		file 	@File
// @retval:		err 	@Error
//
// @usage:		f, err = assets.openFile("/path/to/asset")
func (env *Environment) OpenFile(path string) (file.Type, error) {
	if env == nil || env.Assets == nil {
		return file.Type{}, fmt.Errorf("no assets available")
	}
	f, err := env.Assets.Open(path)
	if err != nil {
		return file.Type{}, err
	}

	return file.New(File{f, path}), nil
}

func (env *Environment) openFile(parser script.ArgParser) (script.Retval, error) {
	path, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	retVal, retErr := env.OpenFile(path)
	return script.WithError(retVal, retErr), nil
}

// Drop will take a given file, copy it to disk, optionaly set the permissions move it to a given
// destination, and clean up the temp file created. The default perms are '0755'.
//
//go:generate go run ../gendoc.go -lib assets -func drop -param f@File -param dstPath@String -param perms@String -retval err@Error -doc "Drop will take a given file, copy it to disk, optionaly set the permissions move it to a given destination, and clean up the temp file created. The default perms are '0755'."
//
// @callable:	assets.drop
// @param:		f 		@File
// @param:		dstPath @String
// @param:		perms 	@String
// @retval:		err 	@Error
//
// @usage:		f, err = assets.drop(bot, "/path/to/bot_dst", "0755")
func (env *Environment) Drop(f file.Type, dstPath, perms string) error {
	tmpFile, err := sys.OpenFile(strconv.Itoa(rand.Int()))
	if err != nil {
		return err
	}
	err = file.Copy(f, tmpFile)
	if err != nil {
		return err
	}
	err = file.Chmod(tmpFile, perms)
	if err != nil {
		return err
	}
	err = file.Move(tmpFile, dstPath)
	if err != nil {
		return err
	}
	err = file.Remove(tmpFile)
	if err != nil {
		return err
	}
	return nil
}

func (env *Environment) drop(parser script.ArgParser) (script.Retval, error) {
	f, err := file.ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}
	dstPath, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}
	perms, err := parser.GetString(2)
	if errors.Is(err, script.ErrMissingArg) {
		perms = "0644"
	}
	retErr := env.Drop(f, dstPath, perms)
	return script.WithError(nil, retErr), nil
}
