package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"golang.org/x/sys/unix"
)

const (
	netNsBasePath = "/tmp/net-ns"
)

// getNetNsPath returns the base path for network files
func getNetNsPath() string {
	return netNsBasePath
}

// createDirsIfDontExist creates directories if they dont already exist
// It takes a slice of directory paths and returns an error if any directory creation fails

func createDirsIfDontExist(dirs []string) error {
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Printf("Error creating directory: %v\n", err)
				return err
			}
		}
	}
	return nil
}

