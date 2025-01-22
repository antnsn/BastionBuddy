//go:build windows
// +build windows

package utils

import (
	"fmt"
	"os/exec"
	"syscall"
)

// GetSysProcAttr returns the system process attributes needed for proper process management
func GetSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		HideWindow: true,
	}
}

// KillProcessGroup kills a process and its entire process group
func KillProcessGroup(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return fmt.Errorf("no process to kill")
	}
	return cmd.Process.Kill()
}

// PrepareAzureCommand prepares an Azure CLI command with the proper environment
func PrepareAzureCommand(args ...string) *exec.Cmd {
	cmd := exec.Command("az.cmd", args...)
	return cmd
}
