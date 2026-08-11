package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStartProxyUsesAvailablePort(t *testing.T) {
	app := &App{}
	proxyPort, err := app.StartProxy(60000)
	if err != nil {
		t.Fatal(err)
	}
	if proxyPort < 1 || proxyPort > 65535 || proxyPort == 60000 {
		t.Fatalf("unexpected proxy port %d", proxyPort)
	}
	if err := app.StopProxy(60000); err != nil {
		t.Fatal(err)
	}
}

func TestProxyForwardsTraffic(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer backend.Close()

	port := backend.Listener.Addr().(*net.TCPAddr).Port
	app := &App{}
	proxyPort, err := app.StartProxy(port)
	if err != nil {
		t.Fatal(err)
	}
	defer app.StopProxy(port)

	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d", proxyPort))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "ok" {
		t.Fatalf("unexpected response %q", body)
	}
}

func TestCaptureBodyPreservesContentAndLimitsLogs(t *testing.T) {
	body := io.NopCloser(strings.NewReader(strings.Repeat("x", maxLoggedBodyBytes+1)))
	restored, logged := captureBody(body)
	defer restored.Close()
	content, err := io.ReadAll(restored)
	if err != nil {
		t.Fatal(err)
	}
	if len(content) != maxLoggedBodyBytes+1 {
		t.Fatalf("expected %d bytes, got %d", maxLoggedBodyBytes+1, len(content))
	}
	if logged != "[Body too large]" {
		t.Fatalf("unexpected log value %q", logged)
	}
}

func TestHeadersForLogRedactsSecrets(t *testing.T) {
	headers := headersForLog(http.Header{
		"Authorization": {"Bearer secret"},
		"Accept":        {"application/json"},
	})
	if headers["Authorization"] != "[REDACTED]" {
		t.Fatalf("authorization header was not redacted")
	}
	if headers["Accept"] != "application/json" {
		t.Fatalf("unexpected accept header %q", headers["Accept"])
	}
}
