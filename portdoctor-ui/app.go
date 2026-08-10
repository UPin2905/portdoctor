package main

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"runtime"

	"github.com/UPin2905/portdoctor/pkg/port"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ScanPorts returns all listening ports
func (a *App) ScanPorts() ([]*port.PortInfo, error) {
	inspector := port.NewInspector()
	return inspector.ListListening()
}

// InspectPort inspects a specific port
func (a *App) InspectPort(p int) (*port.PortInfo, error) {
	inspector := port.NewInspector()
	return inspector.Inspect(p)
}

// KillPort kills the process using a specific port
func (a *App) KillPort(p int) error {
	inspector := port.NewInspector()
	info, err := inspector.Inspect(p)
	if err != nil {
		return err
	}
	if info.PID <= 0 {
		return fmt.Errorf("no process found on port %d", p)
	}

	if runtime.GOOS == "windows" {
		return exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", info.PID)).Run()
	}
	return exec.Command("kill", "-9", fmt.Sprintf("%d", info.PID)).Run()
}

// FindFreePort finds a free port starting from a base port
func (a *App) FindFreePort(base int) (int, error) {
	for p := base; p <= 65535; p++ {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", p))
		if err == nil {
			ln.Close()
			return p, nil
		}
	}
	return 0, fmt.Errorf("no free ports found")
}
