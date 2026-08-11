package main

import "testing"

func TestRedactEnvironmentValue(t *testing.T) {
	if got := redactEnvironmentValue("API_TOKEN", "secret"); got != "[REDACTED]" {
		t.Fatalf("expected redaction, got %q", got)
	}
	if got := redactEnvironmentValue("PATH", "C:\\Tools"); got != "C:\\Tools" {
		t.Fatalf("unexpected path redaction %q", got)
	}
}
