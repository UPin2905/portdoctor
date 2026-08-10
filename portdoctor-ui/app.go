package main

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"runtime"

	"github.com/UPin2905/portdoctor/pkg/port"
	"github.com/UPin2905/portdoctor/pkg/process"
)

// App struct
type App struct {
	ctx context.Context
}

// UIPortInfo extends port.PortInfo with UI-specific fields
type UIPortInfo struct {
	Port        int    `json:"port"`
	Status      string `json:"status"`
	PID         int    `json:"pid"`
	ProcessName string `json:"processName"`
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

// ScanPorts returns all listening ports with process names
func (a *App) ScanPorts() ([]UIPortInfo, error) {
	inspector := port.NewInspector()
	ports, err := inspector.ListListening()
	if err != nil {
		return nil, err
	}

	var results []UIPortInfo
	for _, p := range ports {
		uiInfo := UIPortInfo{
			Port:   p.Port,
			Status: string(p.Status),
			PID:    p.PID,
		}

		if p.PID > 0 {
			proc, err := process.GetProcess(p.PID)
			if err == nil && proc != nil {
				uiInfo.ProcessName = proc.Name
			}
		}
		
		results = append(results, uiInfo)
	}

	return results, nil
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
