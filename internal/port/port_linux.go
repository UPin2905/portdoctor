//go:build linux

package port

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// findPIDForPort trả về PID của process đang listen trên port, 0 nếu không tìm được.
func findPIDForPort(port int) int {
	for _, procFile := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		entries, err := parseProcNetTCP(procFile)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.Port == port && e.PID > 0 {
				return e.PID
			}
		}
	}
	return 0
}

func listListeningPlatform() ([]*PortInfo, error) {
	var results []*PortInfo
	seen := make(map[int]bool)

	for _, procFile := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		entries, err := parseProcNetTCP(procFile)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if seen[e.Port] {
				continue
			}
			seen[e.Port] = true
			results = append(results, e)
		}
	}
	return results, nil
}

func parseProcNetTCP(path string) ([]*PortInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var results []*PortInfo
	scanner := bufio.NewScanner(f)
	scanner.Scan() // skip header

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		// state 0A = LISTEN
		if fields[3] != "0A" {
			continue
		}
		localHex := fields[1]
		portNum, err := parseHexAddr(localHex)
		if err != nil || portNum < 1 {
			continue
		}
		pid := findPIDByInode(fields[9])

		addr := fmt.Sprintf("0.0.0.0:%d", portNum)
		results = append(results, &PortInfo{
			Port:     portNum,
			Protocol: "tcp",
			Address:  addr,
			PID:      pid,
			Status:   StatusOccupied,
		})
	}
	return results, scanner.Err()
}

func parseHexAddr(hexAddr string) (int, error) {
	parts := strings.Split(hexAddr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid hex addr: %s", hexAddr)
	}
	b, err := hex.DecodeString(parts[1])
	if err != nil || len(b) < 2 {
		return 0, err
	}
	// little-endian
	port := int(b[1])<<8 | int(b[0])
	return port, nil
}

func findPIDByInode(inode string) int {
	entries, err := filepath.Glob("/proc/*/fd/*")
	if err != nil {
		return 0
	}
	target := "socket:[" + inode + "]"
	for _, entry := range entries {
		link, err := os.Readlink(entry)
		if err != nil || link != target {
			continue
		}
		parts := strings.Split(entry, "/")
		if len(parts) >= 3 {
			pid, _ := strconv.Atoi(parts[2])
			return pid
		}
	}
	return 0
}
