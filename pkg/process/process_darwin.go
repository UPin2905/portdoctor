//go:build darwin

package process

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getPlatformProcess(pid int) (*ProcessInfo, error) {
	info := &ProcessInfo{PID: pid}

	// ps -p <pid> -o pid,ppid,comm,start,command
	out, err := exec.Command("ps", "-p", strconv.Itoa(pid),
		"-o", "pid=,ppid=,comm=,lstart=,command=").Output()
	if err != nil {
		return info, fmt.Errorf("ps failed: %w", err)
	}

	line := strings.TrimSpace(string(out))
	if line == "" {
		return info, fmt.Errorf("process %d not found", pid)
	}

	fields := strings.Fields(line)
	if len(fields) >= 1 {
		if ppid, err := strconv.Atoi(fields[1]); err == nil {
			info.ParentPID = ppid
		}
	}
	if len(fields) >= 3 {
		info.Name = fields[2]
	}
	// command bắt đầu từ field 8 (sau lstart 4 fields)
	if len(fields) >= 8 {
		info.CommandLine = strings.Join(fields[7:], " ")
	}

	// Executable qua lsof -p
	if exeOut, err := exec.Command("lsof", "-p", strconv.Itoa(pid), "-Fn").Output(); err == nil {
		for _, l := range strings.Split(string(exeOut), "\n") {
			if strings.HasPrefix(l, "n/") && strings.Contains(l, info.Name) {
				info.Executable = strings.TrimPrefix(l, "n")
				break
			}
		}
	}

	// WorkingDir qua lsof cwd
	if cwdOut, err := exec.Command("lsof", "-p", strconv.Itoa(pid), "-a", "-d", "cwd", "-Fn").Output(); err == nil {
		for _, l := range strings.Split(string(cwdOut), "\n") {
			if strings.HasPrefix(l, "n") {
				info.WorkingDir = strings.TrimPrefix(l, "n")
				break
			}
		}
	}

	// Parent name
	if info.ParentPID > 0 {
		if pOut, err := exec.Command("ps", "-p", strconv.Itoa(info.ParentPID), "-o", "comm=").Output(); err == nil {
			info.ParentName = strings.TrimSpace(string(pOut))
		}
	}

	info.StartTime = time.Now() // placeholder

	return info, nil
}
