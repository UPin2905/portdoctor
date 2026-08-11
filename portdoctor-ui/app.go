package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os/exec"
	"regexp"
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
	sharedTunnels map[int]sharedTunnel
	tunnelMutex   sync.Mutex
	ruleEngine    *RuleEngine
}

type sharedTunnel struct {
	command *exec.Cmd
	url     string
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
		sharedTunnels: make(map[int]sharedTunnel),
	}
	a.ruleEngine = NewRuleEngine(a)
	return a
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(context.Context) {
	if a.ruleEngine != nil {
		a.ruleEngine.Stop()
	}

	a.tunnelMutex.Lock()
	defer a.tunnelMutex.Unlock()

	for portNum, tunnel := range a.sharedTunnels {
		if tunnel.command.Process != nil {
			tunnel.command.Process.Kill()
		}
		delete(a.sharedTunnels, portNum)
	}
	stopAllProxies()
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
		if tunnel, exists := a.sharedTunnels[p.Port]; exists {
			uiInfo.SharedURL = tunnel.url
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
	return a.killPort(p, true)
}

func (a *App) killPort(p int, respectProtection bool) error {
	if respectProtection && a.ruleEngine != nil {
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
	if err := port.ValidatePort(portNum); err != nil {
		return "", err
	}

	info, err := port.NewInspector().Inspect(portNum)
	if err != nil {
		return "", err
	}
	if info.PID <= 0 {
		return "", fmt.Errorf("no process found on port %d", portNum)
	}

	a.tunnelMutex.Lock()
	defer a.tunnelMutex.Unlock()

	if _, exists := a.sharedTunnels[portNum]; exists {
		return "", fmt.Errorf("port %d is already being shared", portNum)
	}

	cmd := exec.Command("ssh", "-R", fmt.Sprintf("80:localhost:%d", portNum), "nokey@localhost.run", "-T", "-o", "StrictHostKeyChecking=accept-new")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	lines := make(chan string, 32)
	readOutput := func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
	}
	go readOutput(stdout)
	go readOutput(stderr)
	waitCh := make(chan error, 1)
	go func() {
		waitCh <- cmd.Wait()
	}()

	urlPattern := regexp.MustCompile(`https://[^\s]+`)
	timeout := time.NewTimer(15 * time.Second)
	defer timeout.Stop()
	for {
		select {
		case line := <-lines:
			url := urlPattern.FindString(line)
			if url != "" {
				a.sharedTunnels[portNum] = sharedTunnel{command: cmd, url: url}
				return url, nil
			}
		case err := <-waitCh:
			if err != nil {
				return "", fmt.Errorf("tunnel exited before providing a URL: %w", err)
			}
			return "", fmt.Errorf("tunnel exited before providing a URL")
		case <-timeout.C:
			cmd.Process.Kill()
			return "", fmt.Errorf("timeout waiting for tunnel URL")
		}
	}
}

// StopSharePort stops an active tunnel for a port
func (a *App) StopSharePort(portNum int) error {
	a.tunnelMutex.Lock()
	defer a.tunnelMutex.Unlock()

	tunnel, exists := a.sharedTunnels[portNum]
	if !exists {
		return fmt.Errorf("port %d is not shared", portNum)
	}

	if tunnel.command.Process != nil {
		tunnel.command.Process.Kill()
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
			details.EnvVars[parts[0]] = redactEnvironmentValue(parts[0], parts[1])
		}
	}

	return details, nil
}

func redactEnvironmentValue(key, value string) string {
	normalized := strings.ToUpper(key)
	for _, marker := range []string{"PASSWORD", "SECRET", "TOKEN", "API_KEY", "APIKEY", "CREDENTIAL", "PRIVATE_KEY"} {
		if strings.Contains(normalized, marker) {
			return "[REDACTED]"
		}
	}
	return value
}
