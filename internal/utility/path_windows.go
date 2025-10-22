//go:build windows

package utility

import (
	"path/filepath"
)

func DefaultDir() string {
	return filepath.Join("C:\\", "ProgramData", ServiceName)
}
