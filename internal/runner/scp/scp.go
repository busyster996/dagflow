package scp

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/povsister/scp"
	"golang.org/x/crypto/ssh"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/busyster996/dagflow/internal/common"
	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/internal/utility"
	"github.com/busyster996/dagflow/pkg/logx"
)

type sScp struct {
	storage   storage.IStep
	workspace string

	sshClientConfig *ssh.ClientConfig
	scpClient       *scp.Client

	Method     string        `json:"method"`     // 认证模式
	Username   string        `json:"username"`   // 用户名
	Password   string        `json:"password"`   // 密码
	Host       string        `json:"host"`       // 主机
	Port       int           `json:"port"`       // 端口
	PrivateKey string        `json:"privateKey"` // 私钥
	KeyPass    string        `json:"keyPass"`    // 私钥密码
	Timeout    time.Duration `json:"timeout"`    // 超时时间
	Kind       string        `json:"kind"`       // 类型[文件,目录]
	Direction  string        `json:"direction"`  // 方向[上传,下载]
	Source     string        `json:"source"`     // 源文件/目录
	Target     string        `json:"target"`     // 目标文件/目录
}

func (s *sScp) Run(ctx context.Context) (exit common.ExecCode, err error) {
	defer func() {
		r := recover()
		if r != nil {
			logx.Errorln(string(debug.Stack()), r)
			err = fmt.Errorf("panic: %s", r)
		}
	}()
	err = s.buildConfig()
	if err != nil {
		return common.ExecCodeSystemErr, err
	}
	s.scpClient, err = scp.NewClient(net.JoinHostPort(s.Host, fmt.Sprintf("%d", s.Port)), s.sshClientConfig, &scp.ClientOption{})
	if err != nil {
		logx.Errorf("Failed to connect to the remote server: %s", err)
		return common.ExecCodeSystemErr, err
	}
	if err = s.copy(ctx); err != nil {
		logx.Errorf("Failed to copy: %s", err)
		return common.ExecCodeSystemErr, err
	}
	return common.ExecCodeSuccess, nil
}

func (s *sScp) copy(ctx context.Context) error {
	if s.Kind != "file" {
		fo := &scp.DirTransferOption{
			Context:      ctx,
			Timeout:      s.Timeout,
			PreserveProp: true,
		}
		if s.Direction == "upload" {
			return s.scpClient.CopyDirToRemote(s.Source, s.Target, fo)
		}
		return s.scpClient.CopyDirFromRemote(s.Source, s.Target, fo)
	}
	fo := &scp.FileTransferOption{
		Context:      ctx,
		Timeout:      s.Timeout,
		PreserveProp: true,
	}
	if s.Direction == "upload" {
		return s.scpClient.CopyFileToRemote(s.Source, s.Target, fo)
	}
	return s.scpClient.CopyFileFromRemote(s.Source, s.Target, fo)
}

func (s *sScp) Clear() error {
	if s.scpClient != nil {
		return s.scpClient.Close()
	}
	return nil
}

func (s *sScp) buildConfig() error {
	content, err := s.storage.Content()
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal([]byte(content), s); err != nil {
		return err
	}
	if s.Host == "" {
		return fmt.Errorf("host is empty")
	}
	if s.Port == 0 {
		s.Port = 22
	}
	if s.Kind == "" {
		s.Kind = "file"
	}
	if s.Direction == "" {
		s.Direction = "upload"
	}
	if s.Method == "" {
		s.Method = "password"
	}
	// 判断源路径是否为绝对路径
	s.Source = filepath.Clean(s.Source)
	if !filepath.IsAbs(s.Source) {
		s.Source = filepath.Join(s.workspace, s.Source)
	}

	s.sshClientConfig, err = utility.GetSSHClientConfig(
		s.Method,
		s.Username,
		s.Password,
		s.PrivateKey,
		s.KeyPass,
	)
	if err != nil {
		return err
	}
	return nil
}
