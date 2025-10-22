package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/busyster996/dagflow/internal/common"
	"github.com/busyster996/dagflow/internal/storage"
)

type sTouch struct {
	storage   storage.IStep
	workspace string

	Path      string `json:"path"`
	Overwrite bool   `json:"overwrite"`
	Content   string `json:"content"`
}

func (t *sTouch) Run(ctx context.Context) (exit common.ExecCode, err error) {
	content, err := t.storage.Content()
	if err != nil {
		return common.ExecCodeSystemErr, err
	}
	if err = yaml.Unmarshal([]byte(content), t); err != nil {
		return common.ExecCodeSystemErr, err
	}
	t.Path = filepath.Clean(t.Path)
	if t.Path == "" {
		return common.ExecCodeSystemErr, fmt.Errorf("path is empty")
	}
	if t.Overwrite {
		t.storage.Log().Writef("overwrite %s", t.Path)
		err = os.WriteFile(filepath.Join(t.workspace, t.Path), []byte(t.Content), os.ModePerm)
	} else {
		t.storage.Log().Writef("create or append %s", t.Path)
		var file *os.File
		file, err = os.OpenFile(filepath.Join(t.workspace, t.Path), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		defer func() {
			_ = file.Close()
		}()
		if err != nil {
			return common.ExecCodeSystemErr, err
		}
		_, err = file.WriteString(t.Content)
	}
	if err != nil {
		return common.ExecCodeSystemErr, err
	}
	return common.ExecCodeSuccess, nil
}

func (t *sTouch) Clear() error {
	return nil
}
