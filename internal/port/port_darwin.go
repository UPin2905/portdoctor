//go:build darwin

package port

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

func listListeningPlatform() ([]*PortInfo, error) {
	out, err := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-n", "-P").Output()
	if err != nil {
		return nil, fmt.Errorf("lsof failed: %w", err)
	}
	return parseLsofDarwin(string(out)), nil
}

func parseLsofDarwin(output string) []*PortInfo {
	var results []*PortInfo
	seen := make(map[int]bool)

	for i, line := range strings.Split(output, "\n") {
		if i == 0 {
			continue // skip header
		}
		fields := strings.Fields(line)
		// COMMAND PID USER ... NAME
		if len(fields) < 9 {
			continue
		}
		nameField := fields[8] // e.g. *:3000 or 127.0.0.1:3000
		_, portStr, err := net.SplitHostPort(nameField)
		if err != nil {
			// Try replacing * with 0.0.0.0
			nameField = strings.Replace(nameField, "*:", "0.0.0.0:", 1)
			_, portStr, err = net.SplitHostPort(nameField)
			if err != nil {
				continue
			}
		}
		portNum, err := strconv.Atoi(portStr)
		if err != nil || portNum < 1 {
			continue
		}
		if seen[portNum] {
			continue
		}
		seen[portNum] = true

		pid, _ := strconv.Atoi(fields[1])
		results = append(results, &PortInfo{
			Port:     portNum,
			Protocol: "tcp",
			Address:  nameField,
			PID:      pid,
			Status:   StatusOccupied,
		})
	}
	return results
}
