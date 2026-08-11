package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/UPin2905/portdoctor/pkg/port"
)

type PortRule struct {
	Port           int    `json:"port"`
	Protected      bool   `json:"protected"`
	AllowedProcess string `json:"allowedProcess"`
	AutoHealCmd    string `json:"autoHealCmd"`
	AutoHealDir    string `json:"autoHealDir"`
}

type RuleEngine struct {
	app           *App
	rules         map[int]PortRule
	mutex         sync.Mutex
	rulesFilePath string
	healCooldowns map[int]time.Time
	stop          chan struct{}
	stopOnce      sync.Once
}

func NewRuleEngine(app *App) *RuleEngine {
	appData, _ := os.UserConfigDir()
	rulesDir := filepath.Join(appData, "PortDoctor")
	os.MkdirAll(rulesDir, 0755)

	re := &RuleEngine{
		app:           app,
		rules:         make(map[int]PortRule),
		rulesFilePath: filepath.Join(rulesDir, "rules.json"),
		healCooldowns: make(map[int]time.Time),
		stop:          make(chan struct{}),
	}
	re.LoadRules()
	go re.StartWatcher()
	return re
}

func (re *RuleEngine) LoadRules() {
	re.mutex.Lock()
	defer re.mutex.Unlock()

	data, err := os.ReadFile(re.rulesFilePath)
	if err == nil {
		json.Unmarshal(data, &re.rules)
	}
}

func (re *RuleEngine) SaveRules() error {
	data, err := json.MarshalIndent(re.rules, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(re.rulesFilePath, data, 0644)
}

func (a *App) SaveRule(rule PortRule) error {
	if a.ruleEngine == nil {
		return fmt.Errorf("rule engine not initialized")
	}
	if err := port.ValidatePort(rule.Port); err != nil {
		return err
	}
	a.ruleEngine.mutex.Lock()
	defer a.ruleEngine.mutex.Unlock()

	a.ruleEngine.rules[rule.Port] = rule
	return a.ruleEngine.SaveRules()
}

func (a *App) DeleteRule(portNum int) error {
	if a.ruleEngine == nil {
		return fmt.Errorf("rule engine not initialized")
	}
	a.ruleEngine.mutex.Lock()
	defer a.ruleEngine.mutex.Unlock()

	delete(a.ruleEngine.rules, portNum)
	return a.ruleEngine.SaveRules()
}

func (a *App) GetRules() map[int]PortRule {
	if a.ruleEngine == nil {
		return make(map[int]PortRule)
	}
	a.ruleEngine.mutex.Lock()
	defer a.ruleEngine.mutex.Unlock()

	// Create a copy to return
	rulesCopy := make(map[int]PortRule)
	for k, v := range a.ruleEngine.rules {
		rulesCopy[k] = v
	}
	return rulesCopy
}

func (re *RuleEngine) StartWatcher() {
	inspector := port.NewInspector()
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-re.stop:
			return
		case <-ticker.C:
		}

		re.mutex.Lock()
		rulesCopy := make(map[int]PortRule)
		for k, v := range re.rules {
			rulesCopy[k] = v
		}
		re.mutex.Unlock()

		if len(rulesCopy) == 0 {
			continue
		}

		ports, err := inspector.ListListening()
		if err != nil {
			continue
		}

		activePorts := make(map[int]*port.PortInfo)
		for i := range ports {
			activePorts[ports[i].Port] = ports[i]
		}

		for portNum, rule := range rulesCopy {
			pInfo, isActive := activePorts[portNum]

			if isActive && pInfo.PID > 0 {
				// Port is active
				if rule.AllowedProcess != "" {
					// Check process name
					details, err := re.app.GetProcessDetails(pInfo.PID)
					if err == nil && details != nil {
						if !strings.EqualFold(details.Name, rule.AllowedProcess) && rule.AllowedProcess != "*" {
							// Unauthorized process! Kill it!
							fmt.Printf("RuleEngine: Killing unauthorized process %s on port %d\n", details.Name, portNum)
							re.app.killPort(portNum, false)
						}
					}
				}
			} else if !isActive {
				// Port is dead
				if rule.AutoHealCmd != "" {
					// Check cooldown
					re.mutex.Lock()
					lastHeal, exists := re.healCooldowns[portNum]
					cooldownPassed := !exists || time.Since(lastHeal) > 10*time.Second
					re.mutex.Unlock()

					if cooldownPassed {
						fmt.Printf("RuleEngine: Auto-healing port %d with cmd: %s\n", portNum, rule.AutoHealCmd)
						parts := strings.Fields(rule.AutoHealCmd)
						if len(parts) > 0 {
							cmd := exec.Command(parts[0], parts[1:]...)
							if rule.AutoHealDir != "" {
								cmd.Dir = rule.AutoHealDir
							}
							if err := cmd.Start(); err == nil {
								re.mutex.Lock()
								re.healCooldowns[portNum] = time.Now()
								re.mutex.Unlock()
							}
						}
					}
				}
			}
		}
	}
}

func (re *RuleEngine) Stop() {
	re.stopOnce.Do(func() {
		close(re.stop)
	})
}
