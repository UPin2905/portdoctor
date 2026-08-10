package port

import (
	"net"
	"testing"
)

func TestValidatePort(t *testing.T) {
	valid := []int{1, 80, 3000, 8080, 65535}
	for _, p := range valid {
		if err := ValidatePort(p); err != nil {
			t.Errorf("ValidatePort(%d) unexpected error: %v", p, err)
		}
	}

	invalid := []int{0, -1, 65536, 99999}
	for _, p := range invalid {
		if err := ValidatePort(p); err == nil {
			t.Errorf("ValidatePort(%d) expected error, got nil", p)
		}
	}
}

func TestInspectOccupiedPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not start test listener: %v", err)
	}
	defer ln.Close()

	occupiedPort := ln.Addr().(*net.TCPAddr).Port
	inspector := NewInspector()

	info, err := inspector.Inspect(occupiedPort)
	if err != nil {
		t.Fatalf("Inspect error: %v", err)
	}
	if info.Status != StatusOccupied {
		t.Errorf("expected OCCUPIED, got %s", info.Status)
	}
}

func TestInspectFreePort(t *testing.T) {
	// Lấy port trống từ OS, đóng lại, rồi inspect
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not allocate port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	inspector := NewInspector()
	info, err := inspector.Inspect(port)
	if err != nil {
		t.Fatalf("Inspect error: %v", err)
	}
	if info.Status != StatusFree {
		t.Errorf("expected FREE, got %s", info.Status)
	}
}

func TestInspectInvalidPort(t *testing.T) {
	inspector := NewInspector()
	_, err := inspector.Inspect(0)
	if err == nil {
		t.Error("expected error for port 0, got nil")
	}
	_, err = inspector.Inspect(65536)
	if err == nil {
		t.Error("expected error for port 65536, got nil")
	}
}
