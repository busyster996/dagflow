//go:build windows

package utils

import (
	"fmt"
	"path/filepath"
)

func DefaultDir() string {
	return filepath.Join("C:\\", "ProgramData", ServiceName)
}

func PipeName() string {
	return fmt.Sprintf("npipe:////./pipe/%s", ServiceName)
}
