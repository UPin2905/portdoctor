package framework

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Detect trả về tên framework dựa trên command line và working directory.
func Detect(runtime, cmdLine, workDir string) string {
	cmd := strings.ToLower(cmdLine)

	switch runtime {
	case "Node.js":
		return detectNode(cmd, workDir)
	case "Python":
		return detectPython(cmd)
	}
	return ""
}

func detectNode(cmd, workDir string) string {
	switch {
	case strings.Contains(cmd, "next") && strings.Contains(cmd, "dev"):
		return "Next.js"
	case strings.Contains(cmd, "next"):
		return "Next.js"
	case strings.Contains(cmd, "vite"):
		return "Vite"
	case strings.Contains(cmd, "nuxt"):
		return "Nuxt"
	case strings.Contains(cmd, "nest"):
		return "NestJS"
	case strings.Contains(cmd, "react-scripts"):
		return "Create React App"
	case strings.Contains(cmd, "express"):
		return "Express"
	}

	// Fallback: đọc package.json
	if workDir != "" {
		if fw := detectFromPackageJSON(workDir); fw != "" {
			return fw
		}
	}
	return ""
}

func detectPython(cmd string) string {
	switch {
	case strings.Contains(cmd, "uvicorn"):
		return "Uvicorn / FastAPI"
	case strings.Contains(cmd, "manage.py"):
		return "Django"
	case strings.Contains(cmd, "flask"):
		return "Flask"
	case strings.Contains(cmd, "gunicorn"):
		return "Gunicorn"
	}
	return ""
}

func detectFromPackageJSON(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return ""
	}
	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return ""
	}
	all := make(map[string]bool)
	for k := range pkg.Dependencies {
		all[k] = true
	}
	for k := range pkg.DevDependencies {
		all[k] = true
	}

	switch {
	case all["next"]:
		return "Next.js"
	case all["vite"]:
		return "Vite"
	case all["nuxt"] || all["nuxt3"]:
		return "Nuxt"
	case all["@nestjs/core"]:
		return "NestJS"
	case all["react-scripts"]:
		return "Create React App"
	case all["express"]:
		return "Express"
	case all["svelte"]:
		return "SvelteKit"
	case all["astro"]:
		return "Astro"
	}
	return ""
}
