package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/UPin2905/portdoctor/pkg/port"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// HTTPRequestLog represents a logged HTTP request and its response
type HTTPRequestLog struct {
	ID         string            `json:"id"`
	Method     string            `json:"method"`
	URL        string            `json:"url"`
	ReqHeaders map[string]string `json:"reqHeaders"`
	ReqBody    string            `json:"reqBody"`

	Status     int               `json:"status"`
	ResHeaders map[string]string `json:"resHeaders"`
	ResBody    string            `json:"resBody"`
	LatencyMs  int64             `json:"latencyMs"`
}

var (
	activeProxies map[int]*http.Server
	proxyMutex    sync.Mutex
)

const maxLoggedBodyBytes = 50 * 1024

type preservedBody struct {
	io.Reader
	closer io.Closer
}

func (b *preservedBody) Close() error {
	return b.closer.Close()
}

func init() {
	activeProxies = make(map[int]*http.Server)
}

func captureBody(body io.ReadCloser) (io.ReadCloser, string) {
	prefix, err := io.ReadAll(io.LimitReader(body, maxLoggedBodyBytes+1))
	restored := &preservedBody{
		Reader: io.MultiReader(bytes.NewReader(prefix), body),
		closer: body,
	}
	if err != nil {
		return restored, "[Unable to read body]"
	}
	if len(prefix) > maxLoggedBodyBytes {
		return restored, "[Body too large]"
	}
	return restored, string(prefix)
}

func headersForLog(headers http.Header) map[string]string {
	result := make(map[string]string, len(headers))
	for key, values := range headers {
		switch strings.ToLower(key) {
		case "authorization", "cookie", "proxy-authorization", "set-cookie", "x-api-key":
			result[key] = "[REDACTED]"
		default:
			result[key] = strings.Join(values, ", ")
		}
	}
	return result
}

func (a *App) emit(event string, data any) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, event, data)
	}
}

// StartProxy starts an HTTP reverse proxy for a target port.
// It returns the assigned proxy port.
func (a *App) StartProxy(targetPort int) (int, error) {
	if err := port.ValidatePort(targetPort); err != nil {
		return 0, err
	}

	proxyMutex.Lock()
	defer proxyMutex.Unlock()

	if _, exists := activeProxies[targetPort]; exists {
		return 0, fmt.Errorf("proxy already running for port %d", targetPort)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	proxyPort := listener.Addr().(*net.TCPAddr).Port

	targetURL, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", targetPort))
	if err != nil {
		listener.Close()
		return 0, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Intercept Request and Response
	proxy.ModifyResponse = func(res *http.Response) error {
		reqID := ""
		if res.Request != nil {
			reqID = res.Request.Header.Get("X-PortDoctor-ReqID")
		}

		resLog := make(map[string]interface{})
		resLog["id"] = reqID
		resLog["status"] = res.StatusCode

		resLog["headers"] = headersForLog(res.Header)

		if res.Body != nil {
			res.Body, resLog["body"] = captureBody(res.Body)
		}

		a.emit(fmt.Sprintf("http-res-%d", targetPort), resLog)
		return nil
	}

	director := proxy.Director
	var reqCounter int
	proxy.Director = func(req *http.Request) {
		director(req)

		proxyMutex.Lock()
		reqCounter++
		reqID := fmt.Sprintf("req-%d", reqCounter)
		proxyMutex.Unlock()

		req.Header.Set("X-PortDoctor-ReqID", reqID)

		reqLog := make(map[string]interface{})
		reqLog["id"] = reqID
		reqLog["method"] = req.Method
		reqLog["url"] = req.URL.Path
		if req.URL.RawQuery != "" {
			reqLog["url"] = req.URL.Path + "?" + req.URL.RawQuery
		}

		reqLog["headers"] = headersForLog(req.Header)

		if req.Body != nil {
			req.Body, reqLog["body"] = captureBody(req.Body)
		}

		a.emit(fmt.Sprintf("http-req-%d", targetPort), reqLog)
	}

	server := &http.Server{
		Addr:              fmt.Sprintf("127.0.0.1:%d", proxyPort),
		Handler:           proxy,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	activeProxies[targetPort] = server

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Proxy server error: %v\n", err)
		}
		proxyMutex.Lock()
		if activeProxies[targetPort] == server {
			delete(activeProxies, targetPort)
		}
		proxyMutex.Unlock()
	}()

	return proxyPort, nil
}

// StopProxy stops the reverse proxy
func (a *App) StopProxy(targetPort int) error {
	proxyMutex.Lock()
	defer proxyMutex.Unlock()

	server, exists := activeProxies[targetPort]
	if !exists {
		return fmt.Errorf("proxy not running for port %d", targetPort)
	}

	err := server.Close()
	delete(activeProxies, targetPort)
	return err
}

func stopAllProxies() {
	proxyMutex.Lock()
	servers := make([]*http.Server, 0, len(activeProxies))
	for port, server := range activeProxies {
		servers = append(servers, server)
		delete(activeProxies, port)
	}
	proxyMutex.Unlock()

	for _, server := range servers {
		server.Close()
	}
}
