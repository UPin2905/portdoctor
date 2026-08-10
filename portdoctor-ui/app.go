package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"path/filepath"

	"github.com/UPin2905/portdoctor/pkg/port"
	// Use our internal process for fallback, but gopsutil for advanced features
	"github.com/UPin2905/portdoctor/pkg/process"
	gopsProcess "github.com/shirou/gopsutil/v3/process"
)

// App struct
type App struct {
	ctx           context.Context
	sharedTunnels map[int]*exec.Cmd
	tunnelMutex   sync.Mutex
	ruleEngine    *RuleEngine
}

// UIPortInfo extends port.PortInfo with UI-specific fields
type UIPortInfo struct {
	Port        int     `json:"port"`
	Status      string  `json:"status"`
	PID         int     `json:"pid"`
	ProcessName string  `json:"processName"`
	Project     string  `json:"project"`
	CPU         float64 `json:"cpu"`
	RAM         uint64  `json:"ram"`
	SharedURL   string  `json:"sharedUrl"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	a := &App{
		sharedTunnels: make(map[int]*exec.Cmd),
	}
	a.ruleEngine = NewRuleEngine(a)
	return a
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ScanPorts returns all listening ports with process names and advanced stats
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

			// Extract Project, CPU, RAM using gopsutil
			if gopsProc, err := gopsProcess.NewProcess(int32(p.PID)); err == nil {
				// CPU and RAM
				if cpu, err := gopsProc.CPUPercent(); err == nil {
					uiInfo.CPU = cpu
				}
				if mem, err := gopsProc.MemoryInfo(); err == nil && mem != nil {
					uiInfo.RAM = mem.RSS
				}

				// Project Mapping
				cwd, _ := gopsProc.Cwd()
				if cwd != "" && !strings.Contains(strings.ToLower(cwd), "system32") {
					uiInfo.Project = filepath.Base(cwd)
				} else {
					cmdline, _ := gopsProc.Cmdline()
					// Fallback to cmdline heuristic for project if needed
					if len(cmdline) > 100 {
						uiInfo.Project = "..."
					}
				}
			}
		}

		a.tunnelMutex.Lock()
		if _, exists := a.sharedTunnels[p.Port]; exists {
			uiInfo.SharedURL = "Shared" // We can just set a flag, or we could track the actual URL in a separate map.
		}
		a.tunnelMutex.Unlock()
		
		results = append(results, uiInfo)
	}

	return results, nil
}

// getProcessNames efficiently fetches names for a list of PIDs
func (a *App) getProcessNames(pids []int) map[int]string {
	names := make(map[int]string)
	
	if runtime.GOOS == "windows" {
		snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
		if err != nil {
			return names
		}
		defer syscall.CloseHandle(snapshot)

		var procEntry syscall.ProcessEntry32
		procEntry.Size = uint32(unsafe.Sizeof(procEntry))

		err = syscall.Process32First(snapshot, &procEntry)
		if err != nil {
			return names
		}

		for {
			var buf []uint16
			for _, v := range procEntry.ExeFile {
				if v == 0 {
					break
				}
				buf = append(buf, v)
			}
			name := syscall.UTF16ToString(buf)
			names[int(procEntry.ProcessID)] = name

			err = syscall.Process32Next(snapshot, &procEntry)
			if err != nil {
				break
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
	if a.ruleEngine != nil {
		a.ruleEngine.mutex.Lock()
		if rule, exists := a.ruleEngine.rules[p]; exists && rule.Protected {
			a.ruleEngine.mutex.Unlock()
			return fmt.Errorf("port %d is protected by a rule", p)
		}
		a.ruleEngine.mutex.Unlock()
	}

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
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s", strings.TrimSpace(string(out)))
		}
		return nil
	}
	
	cmd := exec.Command("kill", "-9", fmt.Sprintf("%d", info.PID))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
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

// SharePort exposes the port to the internet via localhost.run
func (a *App) SharePort(portNum int) (string, error) {
	a.tunnelMutex.Lock()
	defer a.tunnelMutex.Unlock()

	if _, exists := a.sharedTunnels[portNum]; exists {
		return "", fmt.Errorf("port %d is already being shared", portNum)
	}

	cmd := exec.Command("ssh", "-R", fmt.Sprintf("80:localhost:%d", portNum), "nokey@localhost.run", "-T", "-o", "StrictHostKeyChecking=no")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	// Read from stdout to find the assigned URL
	urlChan := make(chan string)
	errChan := make(chan error)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			// Expected format: "ddb2150fcd0a93.lhr.life tunneled with tls termination, https://ddb2150fcd0a93.lhr.life"
			if strings.Contains(line, "tunneled with tls termination") {
				parts := strings.Fields(line)
				if len(parts) > 0 {
					urlChan <- parts[0]
					return
				}
			}
		}
		errChan <- fmt.Errorf("failed to extract URL from tunnel output")
	}()

	select {
	case url := <-urlChan:
		a.sharedTunnels[portNum] = cmd
		return "https://" + url, nil
	case err := <-errChan:
		cmd.Process.Kill()
		return "", err
	case <-time.After(15 * time.Second):
		cmd.Process.Kill()
		return "", fmt.Errorf("timeout waiting for tunnel URL")
	}
}

// StopSharePort stops an active tunnel for a port
func (a *App) StopSharePort(portNum int) error {
	a.tunnelMutex.Lock()
	defer a.tunnelMutex.Unlock()

	cmd, exists := a.sharedTunnels[portNum]
	if !exists {
		return fmt.Errorf("port %d is not shared", portNum)
	}

	if cmd.Process != nil {
		cmd.Process.Kill()
	}
	delete(a.sharedTunnels, portNum)
	return nil
}

// ProcessDetails contains detailed information about a process
type ProcessDetails struct {
	PID      int               `json:"pid"`
	Name     string            `json:"name"`
	Cmdline  []string          `json:"cmdline"`
	EnvVars  map[string]string `json:"envVars"`
	Cwd      string            `json:"cwd"`
	Username string            `json:"username"`
}

// GetProcessDetails fetches detailed info like Env Vars and Cmdline for a PID
func (a *App) GetProcessDetails(pid int) (*ProcessDetails, error) {
	if pid <= 0 {
		return nil, fmt.Errorf("invalid PID")
	}

	p, err := gopsProcess.NewProcess(int32(pid))
	if err != nil {
		return nil, fmt.Errorf("could not find process: %w", err)
	}

	details := &ProcessDetails{
		PID:     pid,
		EnvVars: make(map[string]string),
	}

	details.Name, _ = p.Name()
	details.Cwd, _ = p.Cwd()
	details.Cmdline, _ = p.CmdlineSlice()
	details.Username, _ = p.Username()

	env, _ := p.Environ()
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			details.EnvVars[parts[0]] = parts[1]
		}
	}

	return details, nil
}
