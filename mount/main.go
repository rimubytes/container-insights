package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "syscall"
)

// setupNewMountNamespace creates a new mount namespace and sets up the new root filesystem
// It takes the paths for the new root and the directory to put the old root

func setupNewMountNamespace(newRoot, putOld string) error {
    // Bind mount newRoot to itself
    // This is required because the new_root must be a mount point for pivot_root
    if err := syscall.Mount(newRoot, newRoot, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
        return fmt.Errorf("failed to bind mount new root: %v", err)
    } 

    // Create directory for old root
    // put_old must be at or underneath new_root
    oldRootPath := newRoot + putOld
    if err := syscall.Mkdir(oldRootPath, 0700); != nil {
        return fmt.Errof("failed to create directory for old root: %v", err)
    }
}
