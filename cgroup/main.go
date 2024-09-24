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
			//Weight: &weight, // e.g. (weight in the child cgroup) / (sum of cpu weights in the control groups) => percentage of cpu for this child cgroup processes
			Max:    cgroup2.NewCPUMax(&quota, &period), // e.g. 200000 1000000 meaning processes inside this cgroup can (together) run on the CPU for only 0.2 sec every 1 second
			//Cpus:   "0", // This limits on which CPU cores can the processes inside this cgroup run (NOTE: Also "Mems" needs to be set: https://github.com/containerd/cgroups/blob/fa6f6841ed3d57355acadbc06f1d7ed4d91ac4f7/cgroup2/manager.go#L97!)
			//Mems:   "0", // Memory Node” refers to an on-line node that contains memory. 
		},
		Memory: &cgroup2.Memory{
			Max:  pointerInt64(629145600), // ~629MB // If a cgroup's memory usage reaches this limit and can't be reduced, the system OOM killer is invoked on the cgroup. 
			Swap: pointerInt64(314572800), // Swap usage in bytes
			High: pointerInt64(524288000), // memory usage throttle limit. If a cgroup's memory use goes over the high boundary specified here, the cgroup’s processes are throttled and put under heavy reclaim pressure. The default is max, meaning there is no limit.
		},
		IO: &cgroup2.IO{
			Max: []cgroup2.Entry{{
				Major: maj, 
				Minor: min, 
				Type: cgroup2.ReadIOPS, // Limit I/O Read Operations per second for a block device identified as (major, minor) - e.g. "ls -l /dev/sda*"
				Rate: rate, // number of (read) operations per second
			}},
		},
		Pids: &cgroup2.Pids{
			Max: max, // number of processes allowed - The process number controller is used to allow a cgroup hierarchy to stop any new tasks from being fork()’d or clone()’d after a certain limit is reached.
		},
	}
}