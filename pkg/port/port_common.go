package port

import (
	"fmt"
	"net"
)

type platformInspector struct{}

func newPlatformInspector() Inspector {
	return &platformInspector{}
}

// Inspect kiểm tra trạng thái port bằng cách thử listen trên đó.
// Probe cả wildcard và loopback để detect chính xác trên mọi OS.
func (p *platformInspector) Inspect(port int) (*PortInfo, error) {
	if err := ValidatePort(port); err != nil {
		return nil, err
	}

	addrs := []string{
		fmt.Sprintf(":%d", port),          // 0.0.0.0
		fmt.Sprintf("127.0.0.1:%d", port), // loopback
	}

	for _, addr := range addrs {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			pid := findPIDForPort(port)
			return &PortInfo{
				Port:     port,
				Protocol: "tcp",
				Address:  addr,
				PID:      pid,
				Status:   StatusOccupied,
			}, nil
		}
		ln.Close()
	}

	return &PortInfo{
		Port:     port,
		Protocol: "tcp",
		Address:  fmt.Sprintf(":%d", port),
		Status:   StatusFree,
	}, nil
}

// ListListening trả về danh sách các port đang listen.
// Placeholder — sẽ implement đầy đủ ở Phase 7.
func (p *platformInspector) ListListening() ([]*PortInfo, error) {
	return listListeningPlatform()
}
