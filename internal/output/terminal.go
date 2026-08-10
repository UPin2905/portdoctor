package output

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var noColor = os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb"

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

func color(c, s string) string {
	if noColor {
		return s
	}
	return c + s + colorReset
}

// Header in tiêu đề PortDoctor.
func Header() {
	fmt.Println(color(colorBold, "🩺 PortDoctor"))
	fmt.Println()
}

// PortStatus in trạng thái port.
func PortStatus(port int, status string) {
	var c string
	switch status {
	case "FREE":
		c = colorGreen
	case "OCCUPIED":
		c = colorRed
	default:
		c = colorYellow
	}
	fmt.Printf("Port %d is %s\n\n", port, color(c+colorBold, status))
}

// Section in tiêu đề section.
func Section(title string) {
	fmt.Println(color(colorCyan+colorBold, title))
}

// Field in một dòng field: "  Key    Value"
func Field(key, value string) {
	if value == "" {
		return
	}
	fmt.Printf("  %-12s %s\n", key, value)
}

// Diagnosis in phần chẩn đoán.
func Diagnosis(text string) {
	Section("Diagnosis")
	fmt.Printf("  %s\n", text)
	fmt.Println()
}

// SuggestedActions in các lệnh gợi ý.
func SuggestedActions(actions []string) {
	Section("Suggested actions")
	for _, a := range actions {
		fmt.Printf("  %s\n", color(colorCyan, a))
	}
	fmt.Println()
}

// TimeAgo trả về chuỗi "X minutes ago".
func TimeAgo(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	default:
		return t.Format("2006-01-02 15:04")
	}
}

// Error in lỗi theo định dạng thân thiện.
func Error(msg string) {
	fmt.Fprintf(os.Stderr, "%s %s\n", color(colorRed, "Error:"), msg)
}

// TableHeader in header của bảng scan.
func TableHeader() {
	fmt.Printf("%-8s %-16s %-8s %-20s %s\n",
		"PORT", "PROCESS", "PID", "PROJECT", "TYPE")
	fmt.Println(strings.Repeat("─", 64))
}

// Highlight trả về chuỗi được tô màu xanh đậm.
func Highlight(s string) string {
	return color(colorGreen+colorBold, s)
}

// TableRow in một dòng trong bảng scan.
func TableRow(port int, process string, pid int, project, typ string) {
	if project == "" {
		project = "-"
	}
	if typ == "" {
		typ = "-"
	}
	pidStr := fmt.Sprintf("%d", pid)
	if pid == 0 {
		pidStr = "-"
	}
	fmt.Printf("%-8d %-16s %-8s %-20s %s\n", port, process, pidStr, project, typ)
}
