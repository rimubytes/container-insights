// Package main provides functionality for managing network namespaces in Linux.
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	//"strconv"
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	netNsBasePath = "/tmp/net-ns"
)

// getNetNsPath returns the base path for network namespace files.
func getNetNsPath() string {
	return netNsBasePath
}

// createDirsIfDontExist creates directories if they don't already exist.
// It takes a slice of directory paths and returns an error if any directory creation fails.
func createDirsIfDontExist(dirs []string) error {
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, 0755); err != nil {
				log.Printf("Error creating directory: %v\n", err)
				return err
			}
		}
	}
	return nil
}

// setupNewNetworkNamespace creates a new network namespace for the given process ID.
// It performs the following steps:
// 1. Creates necessary directories
// 2. Opens a bind mount file
// 3. Saves the current network namespace
// 4. Creates a new network namespace
// 5. Bind mounts the new namespace to a file
// 6. Sets the process back to the original namespace
func setupNewNetworkNamespace(processID int) {
	if err := createDirsIfDontExist([]string{getNetNsPath()}); err != nil {
		log.Fatalf("Failed to create directories: %v\n", err)
	}

	nsMount := fmt.Sprintf("%s/%d", getNetNsPath(), processID)
	if _, err := syscall.Open(nsMount, syscall.O_RDONLY|syscall.O_CREAT|syscall.O_EXCL, 0644); err != nil {
		log.Fatalf("Unable to open bind mount file: %v\n", err)
	}

	fd, err := syscall.Open("/proc/self/ns/net", syscall.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("Unable to open current network namespace: %v\n", err)
	}
	defer syscall.Close(fd)

	if err := syscall.Unshare(syscall.CLONE_NEWNET); err != nil {
		log.Fatalf("Unshare system call failed: %v\n", err)
	}

	if err := syscall.Mount("/proc/self/ns/net", nsMount, "bind", syscall.MS_BIND, ""); err != nil {
		log.Fatalf("Mount system call failed: %v\n", err)
	}

	if err := unix.Setns(fd, syscall.CLONE_NEWNET); err != nil {
		log.Fatalf("Setns system call failed: %v\n", err)
	}
}

// joinContainerNetworkNamespace sets the current process's network namespace to the one
// specified by the given process ID.
func joinContainerNetworkNamespace(processID int) error {
	nsMount := fmt.Sprintf("%s/%d", getNetNsPath(), processID)
	fd, err := unix.Open(nsMount, unix.O_RDONLY, 0)
	if err != nil {
		log.Printf("Unable to open network namespace file: %v\n", err)
		return err
	}
	defer unix.Close(fd)

	if err := unix.Setns(fd, unix.CLONE_NEWNET); err != nil {
		log.Printf("Setns system call failed: %v\n", err)
		return err
	}
	return nil
}

// getNamespaceInfo retrieves and logs the current network namespace information for a given process ID.
func getNamespaceInfo(processID int) {
	path := fmt.Sprintf("/proc/%d/ns/net", processID)
	out, err := exec.Command("readlink", path).Output()
	if err != nil {
		log.Fatalf("Error reading namespace file: %v\n", err)
	}
	log.Printf("Process is now in the Namespace: %s", string(out))
}

func main() {
	processID := os.Getpid()
	log.Printf("Process ID: %d\n", processID)

	getNamespaceInfo(processID)

	setupNewNetworkNamespace(processID)
	if err := joinContainerNetworkNamespace(processID); err != nil {
		log.Fatalf("Failed to join container network namespace: %v\n", err)
	}

	getNamespaceInfo(processID)
}