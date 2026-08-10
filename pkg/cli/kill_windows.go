//go:build windows

package cli

import (
	"fmt"
	"os/exec"
)

func killWindows(pid int) error {
	return exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", pid)).Run()
}

func killUnix(pid int) error {
	return nil
}
