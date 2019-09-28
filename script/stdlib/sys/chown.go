package sys

import (
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/kcarretto/paragon/script"
)

// Chown uses os.Command to execute the passed string
//
// @param file: A string for execution.
// @param owner: A string representing a user and/or group, separated by a ":" character.
//
// @return (nil, nil) iff success; (nil, err) o/w
func Chown(parser script.ArgParser) (script.Retval, error) {
	file, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	owner, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}

	uid := -1
	gid := -1
	t := strings.Split(owner, ":")

	if t[0] != "" {
		username := t[0]
		userData, err := user.Lookup(username)
		if err != nil {
			return nil, err
		}
		uid32, err := strconv.ParseInt(userData.Uid, 10, 32)
		if err != nil {
			return nil, err
		}
		uid = int(uid32)
	}
	if len(t) > 1 && t[1] != "" {
		group := t[1]
		groupData, err := user.LookupGroup(group)
		if err != nil {
			return nil, err
		}
		gid32, err := strconv.ParseInt(groupData.Gid, 10, 32)
		if err != nil {
			return nil, err
		}
		gid = int(gid32)
	}
	return nil, os.Chown(file, uid, gid)
}
