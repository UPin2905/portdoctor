package process

import "time"

// ProcessInfo chứa thông tin về một process.
type ProcessInfo struct {
	PID         int
	Name        string
	Executable  string
	CommandLine string
	WorkingDir  string
	ParentPID   int
	ParentName  string
	StartTime   time.Time
}

// GetProcess lấy thông tin process từ PID.
// Implementation nằm trong các file process_*.go.
func GetProcess(pid int) (*ProcessInfo, error) {
	return getPlatformProcess(pid)
}
