// +build profile_mem

package main

import (
	"github.com/pkg/profile"
)

func startProfile() interface{ Stop() } {
	return profile.Start(profile.MemProfile)
}
