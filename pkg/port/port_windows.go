//go:build windows

package port

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// findPIDForPort trả về PID của process đang listen trên port, 0 nếu không tìm được.
func findPIDForPort(port int) int {
	cmd := exec.Command("netstat", "-ano")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	target := fmt.Sprintf(":%d", port)
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "TCP") {
			continue
		}
		fields := strings.Fields(line)
		// TCP  0.0.0.0:3306  0.0.0.0:0  LISTENING  1234
		if len(fields) < 5 {
			continue
		}
		if !strings.EqualFold(fields[3], "LISTENING") {
			continue
		}
		localAddr := fields[1]
		if strings.HasSuffix(localAddr, target) {
			pid, _ := strconv.Atoi(fields[4])
			return pid
		}
	}
	return 0
}

func listListeningPlatform() ([]*PortInfo, error) {
	cmd := exec.Command("netstat", "-ano")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("netstat failed: %w", err)
	}
	return parseNetstatWindows(string(out)), nil
}

func parseNetstatWindows(output string) []*PortInfo {
	var results []*PortInfo
	seen := make(map[int]bool)

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "TCP") {
			continue
		}
		fields := strings.Fields(line)
		// TCP  0.0.0.0:3000  0.0.0.0:0  LISTENING  1234
		if len(fields) < 4 {
			continue
		}
		if !strings.EqualFold(fields[3], "LISTENING") {
			continue
		}
		localAddr := fields[1]
		_, portStr, err := net.SplitHostPort(localAddr)
		if err != nil {
			continue
		}
		portNum, err := strconv.Atoi(portStr)
		if err != nil || portNum < 1 {
			continue
		}
		if seen[portNum] {
			continue
		}
		seen[portNum] = true

		pid := 0
		if len(fields) >= 5 {
			pid, _ = strconv.Atoi(fields[4])
		}

		results = append(results, &PortInfo{
			Port:     portNum,
			Protocol: "tcp",
			Address:  localAddr,
			PID:      pid,
			Status:   StatusOccupied,
		})
	}
	return results
}
