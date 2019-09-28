package sys

import (
	"os"

	"github.com/kcarretto/paragon/script"
)

func setBit(n uint32, pos uint) uint32 {
	n |= (1 << pos)
	return n
}

// Chmod uses os.Chmod to change a file's permissions. All optional params are assumed to be false unless passed.
//
// @param file: A string for the path of the file.
// @param ?setUser: A bool for the set user bit.
// @param ?setGroup: A bool for the set group bit.
// @param ?setSticky: A bool for the sticky bit.
// @param ?ownerRead: A bool for the owner read permission.
// @param ?ownerWrite: A bool for the owner write permission. In Windows this is the only bit that matters (set file
//	        to read only iff false; true o/w).
// @param ?ownerExec: A bool for the owner execute permission.
// @param ?groupRead: A bool for the group read permission.
// @param ?groupWrite: A bool for the group write permission.
// @param ?groupExec: A bool for the group execute permission.
// @param ?worldRead: A bool for the world read permission.
// @param ?worldWrite: A bool for the world write permission.
// @param ?worldExec: A bool for the world execute permission.
//
// @return (nil, nil) iff success; (nil, err) o/w
func Chmod(parser script.ArgParser) (script.Retval, error) {
	file, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	var perms uint32
	setUser, _ := parser.GetBoolByName("setUser")
	setGroup, _ := parser.GetBoolByName("setGroup")
	setSticky, _ := parser.GetBoolByName("setSticky")
	ownerRead, _ := parser.GetBoolByName("ownerRead")
	ownerWrite, _ := parser.GetBoolByName("ownerWrite")
	ownerExec, _ := parser.GetBoolByName("ownerExec")
	groupRead, _ := parser.GetBoolByName("groupRead")
	groupWrite, _ := parser.GetBoolByName("groupWrite")
	groupExec, _ := parser.GetBoolByName("groupExec")
	worldRead, _ := parser.GetBoolByName("worldRead")
	worldWrite, _ := parser.GetBoolByName("worldWrite")
	worldExec, _ := parser.GetBoolByName("worldExec")

	// world perms
	if worldExec {
		setBit(perms, 0)
	}
	if worldWrite {
		setBit(perms, 1)
	}
	if worldRead {
		setBit(perms, 2)
	}

	// group perms
	if groupExec {
		setBit(perms, 3)
	}
	if groupWrite {
		setBit(perms, 4)
	}
	if groupRead {
		setBit(perms, 5)
	}

	// owner perms
	if ownerExec {
		setBit(perms, 6)
	}
	if ownerWrite {
		setBit(perms, 7)
	}
	if ownerRead {
		setBit(perms, 8)
	}

	// other perms
	if setSticky {
		setBit(perms, 9)
	}
	if setGroup {
		setBit(perms, 10)
	}
	if setUser {
		setBit(perms, 11)
	}

	err = os.Chmod(file, os.FileMode(perms))
	return nil, err
}
