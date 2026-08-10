package main

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"

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
	pids := make([]int, 0)
	for _, p := range ports {
		if p.PID > 0 {
			pids = append(pids, p.PID)
		}
	}

	names := a.getProcessNames(pids)

	for _, p := range ports {
		uiInfo := UIPortInfo{
			Port:   p.Port,
			Status: string(p.Status),
			PID:    p.PID,
		}

		if p.PID > 0 {
			if name, ok := names[p.PID]; ok {
				uiInfo.ProcessName = name
			}
		}
		
		results = append(results, uiInfo)
	}

	return results, nil
}

// getProcessNames efficiently fetches names for a list of PIDs
func (a *App) getProcessNames(pids []int) map[int]string {
	names := make(map[int]string)
	
	if runtime.GOOS == "windows" {
		cmd := exec.Command("tasklist", "/FO", "CSV", "/NH")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				parts := strings.Split(line, ",")
				if len(parts) >= 2 {
					name := strings.Trim(parts[0], `"`)
					pidStr := strings.Trim(parts[1], `"`)
					if pid, err := strconv.Atoi(pidStr); err == nil {
						names[pid] = name
					}
				}
			}
		}
		return names
	}

	// Fallback for other OSes
	for _, pid := range pids {
		proc, err := process.GetProcess(pid)
		if err == nil && proc != nil {
			names[pid] = proc.Name
		}
	}
	return names
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
		cmd := exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", info.PID))
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		return cmd.Run()
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
