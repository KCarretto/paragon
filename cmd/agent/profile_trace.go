// +build profile_trace

package main

import (
	"github.com/pkg/profile"
)

func startProfile() interface{ Stop() } {
	return profile.Start(profile.TraceProfile)
}
