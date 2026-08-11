package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveRuleValidatesAndPersists(t *testing.T) {
	engine := &RuleEngine{
		rules:         make(map[int]PortRule),
		rulesFilePath: filepath.Join(t.TempDir(), "rules.json"),
	}
	app := &App{ruleEngine: engine}

	if err := app.SaveRule(PortRule{Port: 3000, Protected: true}); err != nil {
		t.Fatal(err)
	}
	if engine.rules[3000].Protected != true {
		t.Fatal("rule was not saved")
	}
	if err := app.SaveRule(PortRule{Port: 0}); err == nil {
		t.Fatal("expected invalid port error")
	}
}

func TestProtectedPortBlocksOnlyManualKill(t *testing.T) {
	app := &App{ruleEngine: &RuleEngine{rules: map[int]PortRule{-1: {Port: -1, Protected: true}}}}
	if err := app.KillPort(-1); err == nil || !strings.Contains(err.Error(), "protected") {
		t.Fatalf("expected protected error, got %v", err)
	}
	if err := app.killPort(-1, false); err != nil && strings.Contains(err.Error(), "protected") {
		t.Fatalf("automated kill should bypass manual protection: %v", err)
	}
}
