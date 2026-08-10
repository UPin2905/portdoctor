//go:build windows

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

	// Lấy Name + Executable bằng tasklist
	if err := fillFromTasklist(info); err != nil {
		// Không critical — tiếp tục với thông tin còn lại
		info.Name = "unknown"
	}

	// Lấy CommandLine + WorkingDir + ParentPID qua WMIC
	fillFromWMIC(info)

	return info, nil
}

func fillFromTasklist(info *ProcessInfo) error {
	out, err := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", info.PID), "/FO", "CSV", "/NH").Output()
	if err != nil {
		return err
	}
	line := strings.TrimSpace(string(out))
	if line == "" || strings.Contains(line, "No tasks") {
		return fmt.Errorf("process %d not found", info.PID)
	}
	// CSV: "name.exe","pid","session","num","mem"
	parts := strings.Split(line, ",")
	if len(parts) >= 1 {
		info.Name = strings.Trim(parts[0], `"`)
	}
	return nil
}

func fillFromWMIC(info *ProcessInfo) {
	// CommandLine
	out, err := exec.Command("wmic", "process", "where",
		fmt.Sprintf("ProcessId=%d", info.PID),
		"get", "CommandLine,ExecutablePath,ParentProcessId,CreationDate", "/FORMAT:CSV").Output()
	if err != nil {
		return
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Node") {
			continue
		}
		// CSV: Node,CommandLine,CreationDate,ExecutablePath,ParentProcessId
		parts := strings.SplitN(line, ",", 5)
		if len(parts) < 5 {
			continue
		}
		info.CommandLine = parts[1]
		info.Executable = parts[3]
		if ppid, err := strconv.Atoi(parts[4]); err == nil {
			info.ParentPID = ppid
		}
		// CreationDate format: 20240101120000.000000+420
		if t, err := parseWMICDate(parts[2]); err == nil {
			info.StartTime = t
		}
	}

	// WorkingDir — lấy qua PowerShell vì WMIC không có trực tiếp
	out, err = exec.Command("powershell", "-NoProfile", "-Command",
		fmt.Sprintf("(Get-Process -Id %d).Path | Split-Path", info.PID)).Output()
	if err == nil {
		info.WorkingDir = strings.TrimSpace(string(out))
	}

	// Parent name
	if info.ParentPID > 0 {
		parentOut, err := exec.Command("tasklist", "/FI",
			fmt.Sprintf("PID eq %d", info.ParentPID), "/FO", "CSV", "/NH").Output()
		if err == nil {
			pLine := strings.TrimSpace(string(parentOut))
			pParts := strings.Split(pLine, ",")
			if len(pParts) >= 1 {
				info.ParentName = strings.Trim(pParts[0], `"`)
			}
		}
	}
}

func parseWMICDate(s string) (time.Time, error) {
	// 20240101120000.000000+420
	if len(s) < 14 {
		return time.Time{}, fmt.Errorf("short date: %s", s)
	}
	return time.Parse("20060102150405", s[:14])
}
