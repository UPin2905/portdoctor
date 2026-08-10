//go:build linux

package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func getPlatformProcess(pid int) (*ProcessInfo, error) {
	info := &ProcessInfo{PID: pid}
	base := fmt.Sprintf("/proc/%d", pid)

	// Name từ /proc/<pid>/comm
	if b, err := os.ReadFile(filepath.Join(base, "comm")); err == nil {
		info.Name = strings.TrimSpace(string(b))
	}

	// Executable từ /proc/<pid>/exe symlink
	if exe, err := os.Readlink(filepath.Join(base, "exe")); err == nil {
		info.Executable = exe
	}

	// CommandLine từ /proc/<pid>/cmdline (null-separated)
	if b, err := os.ReadFile(filepath.Join(base, "cmdline")); err == nil {
		info.CommandLine = strings.ReplaceAll(string(b), "\x00", " ")
		info.CommandLine = strings.TrimSpace(info.CommandLine)
	}

	// WorkingDir từ /proc/<pid>/cwd symlink
	if cwd, err := os.Readlink(filepath.Join(base, "cwd")); err == nil {
		info.WorkingDir = cwd
	}

	// ParentPID + StartTime từ /proc/<pid>/stat
	if b, err := os.ReadFile(filepath.Join(base, "stat")); err == nil {
		fillFromStat(info, string(b))
	}

	// ParentName
	if info.ParentPID > 0 {
		if b, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", info.ParentPID)); err == nil {
			info.ParentName = strings.TrimSpace(string(b))
		}
	}

	return info, nil
}

func fillFromStat(info *ProcessInfo, stat string) {
	// Format: pid (name) state ppid ...
	// Tìm closing paren để skip tên process (có thể chứa spaces)
	end := strings.LastIndex(stat, ")")
	if end < 0 {
		return
	}
	rest := strings.TrimSpace(stat[end+1:])
	fields := strings.Fields(rest)
	// fields[0] = state, fields[1] = ppid, fields[2] = pgrp...
	// fields[19] = starttime (jiffies since boot)
	if len(fields) > 1 {
		if ppid, err := strconv.Atoi(fields[1]); err == nil {
			info.ParentPID = ppid
		}
	}
	if len(fields) > 19 {
		if jiffies, err := strconv.ParseInt(fields[19], 10, 64); err == nil {
			info.StartTime = jiffiesToTime(jiffies)
		}
	}
}

func jiffiesToTime(jiffies int64) time.Time {
	// Đây là approximation — chính xác cần đọc /proc/uptime + boot time
	// Đủ dùng cho MVP
	b, err := os.ReadFile("/proc/stat")
	if err != nil {
		return time.Time{}
	}
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "btime ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if btime, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
					startSec := btime + jiffies/100
					return time.Unix(startSec, 0)
				}
			}
		}
	}
	return time.Time{}
}
