package project

import (
	"os"
	"path/filepath"
)

// ProjectInfo chứa thông tin về dự án.
type ProjectInfo struct {
	Root      string
	Name      string
	Runtime   string
	Framework string
}

// markers là danh sách file/thư mục đánh dấu project root.
var markers = []string{
	".git", "go.mod", "package.json", "pyproject.toml",
	"requirements.txt", "Cargo.toml", "pom.xml", "build.gradle",
	"composer.json", "Gemfile",
}

// Detect tìm project root bắt đầu từ dir, đi lên tối đa maxDepth cấp.
func Detect(dir string) *ProjectInfo {
	if dir == "" {
		return nil
	}

	current := dir
	for i := 0; i < 8; i++ {
		for _, marker := range markers {
			candidate := filepath.Join(current, marker)
			if _, err := os.Stat(candidate); err == nil {
				return &ProjectInfo{
					Root: current,
					Name: filepath.Base(current),
				}
			}
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return nil
}
