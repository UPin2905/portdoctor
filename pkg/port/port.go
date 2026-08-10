package port

import "fmt"

// Status biểu thị trạng thái của một port.
type Status string

const (
	StatusFree     Status = "FREE"
	StatusOccupied Status = "OCCUPIED"
	StatusUnknown  Status = "UNKNOWN"
)

// PortInfo chứa thông tin đầy đủ về một port.
type PortInfo struct {
	Port     int
	Protocol string
	Address  string
	PID      int
	Status   Status
}

// Inspector là interface cho việc kiểm tra port, tách biệt OS-specific code.
type Inspector interface {
	Inspect(port int) (*PortInfo, error)
	ListListening() ([]*PortInfo, error)
}

// NewInspector trả về Inspector phù hợp với OS hiện tại.
// Implementation nằm trong các file port_*.go.
func NewInspector() Inspector {
	return newPlatformInspector()
}

// ValidatePort kiểm tra port có hợp lệ không (1–65535).
func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid port %d: must be between 1 and 65535", port)
	}
	return nil
}
