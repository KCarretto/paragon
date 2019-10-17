// +build !profile_cpu,!profile_mem,!profile_trace

package main

type stopper bool

func (s stopper) Stop() {}

func startProfile() interface{ Stop() } {
	var s stopper
	return s
}
