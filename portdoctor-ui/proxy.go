package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// HTTPRequestLog represents a logged HTTP request and its response
type HTTPRequestLog struct {
	ID         string `json:"id"`
	Method     string `json:"method"`
	URL        string `json:"url"`
	ReqHeaders map[string]string `json:"reqHeaders"`
	ReqBody    string `json:"reqBody"`

	Status     int               `json:"status"`
	ResHeaders map[string]string `json:"resHeaders"`
	ResBody    string            `json:"resBody"`
	LatencyMs  int64             `json:"latencyMs"`
}

var (
	activeProxies map[int]*http.Server
	proxyMutex    sync.Mutex
)

func init() {
	activeProxies = make(map[int]*http.Server)
}

// StartProxy starts an HTTP reverse proxy for a target port.
// It returns the assigned proxy port.
func (a *App) StartProxy(targetPort int) (int, error) {
	proxyMutex.Lock()
	defer proxyMutex.Unlock()

	// Pick a proxy port, e.g., targetPort + 10000
	proxyPort := targetPort + 10000
	if _, exists := activeProxies[targetPort]; exists {
		return proxyPort, nil // Already running
	}

	targetURL, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", targetPort))
	if err != nil {
		return 0, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Intercept Request and Response
	proxy.ModifyResponse = func(res *http.Response) error {
		reqID := res.Request.Header.Get("X-PortDoctor-ReqID")
		
		resLog := make(map[string]interface{})
		resLog["id"] = reqID
		resLog["status"] = res.StatusCode
		
		headers := make(map[string]string)
		for k, v := range res.Header {
			headers[k] = v[0]
		}
		resLog["headers"] = headers

		if res.Body != nil {
			bodyBytes, _ := io.ReadAll(res.Body)
			res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			if len(bodyBytes) < 50000 { 
				resLog["body"] = string(bodyBytes)
			} else {
				resLog["body"] = "[Body too large]"
			}
		}

		runtime.EventsEmit(a.ctx, fmt.Sprintf("http-res-%d", targetPort), resLog)
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
		
		headers := make(map[string]string)
		for k, v := range req.Header {
			headers[k] = v[0]
		}
		reqLog["headers"] = headers

		if req.Body != nil {
			bodyBytes, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			if len(bodyBytes) < 50000 {
				reqLog["body"] = string(bodyBytes)
			}
		}

		runtime.EventsEmit(a.ctx, fmt.Sprintf("http-req-%d", targetPort), reqLog)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", proxyPort),
		Handler: proxy,
	}

	activeProxies[targetPort] = server

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Proxy server error: %v\n", err)
		}
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
