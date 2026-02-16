package treemux

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// normalizePath returns an absolute, cleaned path.
func normalizePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("missing path")
	}
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, strings.TrimPrefix(path, "~"))
	}

	if !filepath.IsAbs(path) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		path = filepath.Join(cwd, path)
	}

	return filepath.Clean(path), nil
}
