package config

import (
	"path/filepath"
)

func normalizePath(path, baseDir string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join(baseDir, path)
	}
	return filepath.Clean(path)
}
