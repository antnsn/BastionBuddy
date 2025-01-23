//go:build !windows
// +build !windows

package utils

import (
	"fmt"
	"os/exec"
	"syscall"
)

// GetSysProcAttr returns the system process attributes needed for proper process management
func GetSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setpgid: true, // Create new process group
	}
}

// KillProcessGroup kills a process and its entire process group
func KillProcessGroup(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return fmt.Errorf("no process to kill")
	}

	// On Unix-like systems, kill the process group
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to kill process group: %v", err)
	}

	return nil
}

// PrepareAzureCommand prepares an Azure CLI command with the proper environment
func PrepareAzureCommand(args ...string) *exec.Cmd {
	cmd := exec.Command("az", args...)
	return cmd
}
