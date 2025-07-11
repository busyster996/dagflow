//go:build !windows

package utils

import (
	"fmt"
	"path/filepath"
)

func DefaultDir() string {
	return filepath.Join("/", "usr", "local", ServiceName)
}

func PipeName() string {
	return fmt.Sprintf("unix:///var/run/%s.sock", ServiceName)
}
