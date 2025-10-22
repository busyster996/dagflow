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

type sMkdir struct {
	storage   storage.IStep
	workspace string

	Path string `json:"path"` // 文件夹路径
}

func (m *sMkdir) Run(ctx context.Context) (exit common.ExecCode, err error) {
	content, err := m.storage.Content()
	if err != nil {
		return common.ExecCodeSystemErr, err
	}
	if err = yaml.Unmarshal([]byte(content), m); err != nil {
		return common.ExecCodeSystemErr, err
	}
	m.Path = filepath.Clean(m.Path)
	if m.Path == "" {
		return common.ExecCodeSystemErr, fmt.Errorf("path is empty")
	}
	m.storage.Log().Writef("mkdir -p %s", m.Path)
	err = os.MkdirAll(filepath.Join(m.workspace, m.Path), os.ModePerm)
	if err != nil {
		return common.ExecCodeSystemErr, err
	}
	return common.ExecCodeSuccess, nil
}

func (m *sMkdir) Clear() error {
	return nil
}
