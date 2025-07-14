package runner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"golang.org/x/net/context"

	"github.com/busyster996/dagflow/internal/storage"
)

var runners = make(map[string]Factory)

type Factory func(storage storage.IStep, subCmd, workspace, scriptDir string) (IRunner, error)

func Register(name string, factory Factory) {
	runners[strings.ToLower(name)] = factory
}

func ListAvailable() []string {
	var names []string
	for name := range runners {
		names = append(names, name)
	}
	return names
}

func Get(name string) (Factory, error) {
	factory, ok := runners[strings.ToLower(name)]
	if !ok {
		// fallback to exec
		factory, ok = runners["exec"]
		if !ok {
			return nil, errors.New("unknown executor type")
		}
	}
	return factory, nil
}

// default runner
func init() {
	// sh/bash/cmd/powershell/python(2/3) runner
	Register("exec", func(storage storage.IStep, subCmd, workspace, scriptDir string) (IRunner, error) {
		var c = &sCmd{
			storage:   storage,
			workspace: workspace,
			shell:     subCmd,
		}
		c.ctx, c.cancel = context.WithCancel(context.Background())
		c.scriptName = filepath.Join(scriptDir, ksuid.New().String())
		c.scriptName = c.scriptName + c.scriptSuffix()
		if err := os.MkdirAll(scriptDir, os.ModePerm); err != nil {
			return nil, err
		}
		content, err := storage.Content()
		if err != nil {
			return nil, err
		}
		if c.shell == "cmd" || c.shell == "powershell" {
			content = c.utf8ToGb2312(content)
		}
		if err = os.WriteFile(c.scriptName, []byte(content), os.ModePerm); err != nil {
			return nil, err
		}
		return c, nil
	})
	// mkdir runner
	Register("mkdir", func(storage storage.IStep, subCmd, workspace, scriptDir string) (IRunner, error) {
		return &sMkdir{
			storage:   storage,
			workspace: workspace,
		}, nil
	})
	// touch runner
	Register("touch", func(storage storage.IStep, subCmd, workspace, scriptDir string) (IRunner, error) {
		return &sTouch{
			storage:   storage,
			workspace: workspace,
		}, nil
	})
}
