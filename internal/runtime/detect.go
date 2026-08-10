package runtime

import (
	"path/filepath"
	"strings"
)

// Detect trả về tên runtime dựa trên process name và command line.
func Detect(processName, cmdLine string) string {
	name := strings.ToLower(filepath.Base(processName))
	// Xóa .exe suffix trên Windows
	name = strings.TrimSuffix(name, ".exe")
	cmd := strings.ToLower(cmdLine)

	switch {
	case name == "node" || name == "nodejs" || strings.Contains(cmd, "node "):
		return "Node.js"
	case name == "python" || name == "python3" || name == "python2":
		return "Python"
	case name == "java":
		return "Java"
	case name == "dotnet" || strings.HasSuffix(name, ".dll"):
		return ".NET"
	case name == "go" || strings.Contains(cmd, "go run"):
		return "Go"
	case name == "ruby":
		return "Ruby"
	case name == "php":
		return "PHP"
	case name == "cargo" || name == "rust":
		return "Rust"
	case name == "docker-proxy" || name == "docker" || name == "containerd":
		return "Docker"
	}
	return ""
}
