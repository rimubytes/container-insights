// Package main provides a program to demonstrate cgroup management using the cgroups/v3 library.
package main

import (
	"log"
	"os/exec"
	"time"

	"github.com/containerd/cgroups/v3/cgroup2"
)

// pointerInt64 returns a pointer to the given int64 value.
// This is useful for cgroup configurations that require pointers.
func pointerInt64(i int64) *int64 {
	return &i
}

// setupCgroupResources creates and returns a cgroup2.Resources struct
// with all the resource limits configured.
func setupCgroupResources() cgroup2.Resources {
	return cgroup2.Resources{
		CPU:    setupCPUResources(),
		Memory: setupMemoryResources(),
		IO:     setupIOResources(),
		Pids:   setupPidsResources(),
	}
}

// setupCPUResources configures and returns CPU-specific resource limits.
func setupCPUResources() *cgroup2.CPU {
	quota := int64(200000)
	period := uint64(1000000)
	return &cgroup2.CPU{
		// Max CPU usage: 200ms per 1000ms (20% of CPU)
		Max: cgroup2.NewCPUMax(&quota, &period),
	}
}

// setupMemoryResources configures and returns memory-specific resource limits.
func setupMemoryResources() *cgroup2.Memory {
	return &cgroup2.Memory{
		Max:  pointerInt64(629145600), // ~629MB max memory usage
		Swap: pointerInt64(314572800), // ~300MB max swap usage
		High: pointerInt64(524288000), // ~500MB memory throttle limit
	}
}

// setupIOResources configures and returns I/O-specific resource limits.
func setupIOResources() *cgroup2.IO {
	return &cgroup2.IO{
		Max: []cgroup2.Entry{{
			Major: 8,
			Minor: 0,
			Type:  cgroup2.ReadIOPS,
			Rate:  120, // 120 read operations per second
		}},
	}
}

// setupPidsResources configures and returns process ID limits.
func setupPidsResources() *cgroup2.Pids {
	return &cgroup2.Pids{
		Max: 1000, // Maximum number of processes
	}
}

// createCgroup creates a new cgroup with the specified resources.
func createCgroup(res *cgroup2.Resources) (*cgroup2.Manager, error) {
	return cgroup2.NewSystemd("/", "my-cgroup-abc.slice", -1, res)
}

// runStressCommand starts a stress command to simulate CPU load.
func runStressCommand() (*exec.Cmd, error) {
	cmd := exec.Command("stress", "-c", "1", "--timeout", "30")
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func main() {
	// Set up cgroup resources
	res := setupCgroupResources()
	
	// Create the cgroup
	m, err := createCgroup(&res)
	if err != nil {
		log.Fatalf("Error creating cgroup: %v\n", err)
	}

	// Get and log the cgroup type
	cgType, err := m.GetType()
	if err != nil {
		log.Fatalf("Error getting cgroup type: %v\n", err)
	}
	log.Println("Cgroup type:", cgType)

	// Run the stress command
	cmd, err := runStressCommand()
	if err != nil {
		log.Fatalf("Error starting stress command: %v\n", err)
	}

	// Get the PID of the stress command (adding 1 because stress spawns a child process)
	pid := cmd.Process.Pid + 1
	log.Printf("PID of the spawned process: %d\n", pid)

	// Add the process to the cgroup
	if err := m.AddProc(uint64(pid)); err != nil {
		log.Fatalf("Error adding process to cgroup: %v\n", err)
	}

	// List processes in the cgroup
	procs, err := m.Procs(false)
	if err != nil {
		log.Fatalf("Error getting processes in cgroup: %v\n", err)
	}
	log.Printf("List of processes inside this cgroup: %v", procs)

	// Freeze the cgroup
	log.Println("Freezing Process")
	if err := m.Freeze(); err != nil {
		log.Fatalf("Error freezing process: %v\n", err)
	}

	// Wait for 15 seconds
	time.Sleep(time.Second * 15)

	// Thaw the cgroup
	log.Println("Thawing Process")
	if err := m.Thaw(); err != nil {
		log.Fatalf("Error thawing process: %v\n", err)
	}

	// Wait for the stress command to finish
	if err := cmd.Wait(); err != nil {
		log.Printf("Error waiting for the command to finish: %v\n", err)
	}

	// Clean up: delete the cgroup
	if err := m.DeleteSystemd(); err != nil {
		log.Fatalf("Error deleting cgroup: %v\n", err)
	}
}