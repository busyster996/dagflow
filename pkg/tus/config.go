package tus

import (
	"fmt"
	"net/url"

	"github.com/busyster996/dagflow/pkg/tus/storage"
	"github.com/busyster996/dagflow/pkg/tus/types"
)

type SConfig struct {
	MaxSize                    int64
	BasePath                   string
	isAbs                      bool
	Store                      storage.IStorage
	Logger                     ILogger
	PreUploadCreateCallback    func(hook types.HookEvent) (types.HTTPResponse, types.FileInfoChanges, error)
	PreFinishResponseCallback  func(hook types.HookEvent) (types.HTTPResponse, error)
	PreUploadTerminateCallback func(hook types.HookEvent) (types.HTTPResponse, error)
}

func (config *SConfig) validate() error {
	if config.Logger == nil {
		return fmt.Errorf("logger is required")
	}

	base := config.BasePath
	uri, err := url.Parse(base)
	if err != nil {
		return err
	}

	// Ensure base path ends with slash to remove logic from absFileURL
	if base != "" && string(base[len(base)-1]) != "/" {
		base += "/"
	}

	// Ensure base path begins with slash if not absolute (starts with scheme)
	if !uri.IsAbs() && len(base) > 0 && string(base[0]) != "/" {
		base = "/" + base
	}

	config.BasePath = base
	config.isAbs = uri.IsAbs()
	return nil
}
