//go:build linux || darwin

package cli

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func killWindows(pid int) error {
	return nil
}

func killUnix(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("process %d not found: %w", pid, err)
	}

	// Graceful: SIGTERM
	if err := p.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("SIGTERM failed: %w", err)
	}

	// Chờ 3 giây
	for i := 0; i < 30; i++ {
		time.Sleep(100 * time.Millisecond)
		if err := p.Signal(syscall.Signal(0)); err != nil {
			// Process đã chết
			return nil
		}
	}

	// Force: SIGKILL
	return p.Signal(syscall.SIGKILL)
}
