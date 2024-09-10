package main

import (
	"log"
	"time"
	"os/exec"
	"github.com/containerd/cgroups/v3/cgroup2"
)

func pointerInt64(int int64) *int64 {
	return &int 
}

func main() {
	var (
		quota int64 = 200000
		period uint = 1000000
		maj int64 = 8
		min int64 = 0
		rate uint64 = 120
		max int64 = 1000
	)

	res := cgroup2.Resources{
		CPU: &cgroup2.CPU{
			Max: cgroup2.NewCPUMax(&quota, &period), //CPU time: 200ms per 1 second
		},
		Memory: &cgroup2.Memory{
			Max: pointerInt64(629145600),
			Swap: pointerInt64(314572800),
			High: pointerInt64(52488000),
		},
		IO: &cgroup2.IO{
			Max: []cgroup2.Entry{{
				Major: mj,
				Minor: min,
				Type: cgroup2.ReadIOPS,
				Rate: rate,
			}},
		},
		Pids: &cgroup2.Pids{
			Max: max,
		},
	}
}