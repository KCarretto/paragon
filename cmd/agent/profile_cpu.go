// +build profile_cpu

package main

import (
	"github.com/pkg/profile"
)

func startProfile() interface{ Stop() } {
	return profile.Start()
}
