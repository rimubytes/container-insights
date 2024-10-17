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
        return fmt.Errorf("failed to create directory for old root: %v", err)
    }

	// Create a new mount namespace
	if err := syscall.Unshare(syscall.CLONE_NEWNS); err != nil {
		return fmt.Errorf("failed to unshare mount namespace: %v", err)
	}

	// Change the root filesystem
	if err := syscall.PivotRoot(newRoot, oldRootPath); err != nil {
		return fmt.Errorf("failed to pivot root: %v", err)
	}

	// Change the current working directory
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("failed to change directory to /: %v", err)
	}

	// Mount essential filesystems
	if err := mountEssentialFS(); err != nil {
		return err
	}

	// Unmount the old root filesystem
	if err := syscall.Unmount(putOld, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("failed to unmount old root filesystem: %v", err)
	}

	return nil
}

// mountEssentialFS mounts essentail filesystems in the new namespace
func mountEssentialFS() error {
    mounts: := [] struct {
        source string
        target string
        fstype string
        flags uintptr
        data string
    }{
        {"/proc", "/proc", 0, ""}
        {"/dev", "/dev", "tmpfs", 0, ""}
    }

    // Create /dev/null
    if _,err := os.Create("/dev/null"); err != nil {
        retun fmt.Errorf("failed to create /dev/null: %v", err)
    }

    return nil
}

// getMountNamespace retrieves the current mount namespace
func getMountNamespace() (string, error) {
    out, err := exec.Command("readlink", "/proc/self/ns/mnt").Output()
    if err != nil {
        return "", fmt.Errorf("error reading namespace file: %v", err)
    }
    return string(out), nil
}

// runShell starts a shell in the new namespace
func runShell() error {
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	processID := os.Getpid()
	log.Printf("Process ID: %d\n", processID)

	oldNS, err := getMountNamespace()
	if err != nil {
		log.Fatalf("Failed to get old mount namespace: %v", err)
	}
	log.Printf("Process is in the old mount Namespace: %s", oldNS)

	newRoot := "new_root"
	putOld := "/old_root"

	if err := setupNewMountNamespace(newRoot, putOld); err != nil {
		log.Fatalf("Failed to setup new mount namespace: %v", err)
	}

	newNS, err := getMountNamespace()
	if err != nil {
		log.Fatalf("Failed to get new mount namespace: %v", err)
	}
	log.Printf("Process is now in the new mount Namespace: %s", newNS)

	log.Println("Opening a shell (/bin/sh) in the new mount namespace")
	log.Println("You can run commands like `mount`, 'lsns', etc.")

	if err := runShell(); err != nil {
		log.Fatalf("Failed to run shell: %v", err)
	}
}