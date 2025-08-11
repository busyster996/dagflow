//go:build !windows

package utils

import (
	"path/filepath"
)

func DefaultDir() string {
	return filepath.Join("/", "usr", "local", ServiceName)
}
