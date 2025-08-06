package sftp

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/internal/worker/common"
	"github.com/busyster996/dagflow/pkg/logx"
)

type sSftp struct {
	storage   storage.IStep
	workspace string

	sshClientConfig *ssh.ClientConfig
	sshClient       *ssh.Client
	sftpClient      *sftp.Client

	Method     string `json:"method"`     // 认证模式
	Username   string `json:"username"`   // 用户名
	Password   string `json:"password"`   // 密码
	Host       string `json:"host"`       // 主机
	Port       int    `json:"port"`       // 端口
	PrivateKey string `json:"privateKey"` // 私钥
	KeyPass    string `json:"keyPass"`    // 私钥密码
	Direction  string `json:"direction"`  // 方向[上传,下载]
	Source     string `json:"source"`     // 源文件/目录
	Target     string `json:"target"`     // 目标文件/目录
}

func (s *sSftp) Run(ctx context.Context) (exit int64, err error) {
	if err = s.buildConfig(); err != nil {
		logx.Errorln(err)
		return common.CodeSystemErr, err
	}
	s.sshClient, err = ssh.Dial("tcp", net.JoinHostPort(s.Host, fmt.Sprintf("%d", s.Port)), s.sshClientConfig)
	if err != nil {
		logx.Errorf("Failed to connect to the remote server: %s", err)
		return common.CodeSystemErr, err
	}
	s.sftpClient, err = sftp.NewClient(s.sshClient)
	if err != nil {
		logx.Errorf("Failed to create sftp client: %s", err)
		return common.CodeSystemErr, err
	}
	if err = s.copy(); err != nil {
		logx.Errorf("Failed to copy: %s", err)
		return common.CodeSystemErr, err
	}
	return common.CodeSuccess, nil
}

func (s *sSftp) copy() error {
	if s.Direction == "upload" {
		return s.uploadFile()
	}
	return s.downloadFile()
}

func (s *sSftp) uploadFile() error {
	lf, err := os.Open(s.Source)
	if err != nil {
		return fmt.Errorf("open local file: %w", err)
	}
	defer lf.Close()

	rf, err := s.sftpClient.Create(s.Target)
	if err != nil {
		return fmt.Errorf("create remote file: %w", err)
	}
	defer rf.Close()

	if _, err = io.Copy(rf, lf); err != nil {
		return fmt.Errorf("copy to remote: %w", err)
	}
	return nil
}

func (s *sSftp) downloadFile() error {
	rf, err := s.sftpClient.Open(s.Source)
	if err != nil {
		return fmt.Errorf("open remote file: %w", err)
	}
	defer rf.Close()

	_ = os.MkdirAll(filepath.Dir(s.Target), os.ModeDir)
	lf, err := os.Create(s.Target)
	if err != nil {
		return fmt.Errorf("create local file: %w", err)
	}
	defer lf.Close()

	if _, err = io.Copy(lf, rf); err != nil {
		return fmt.Errorf("copy to local: %w", err)
	}
	return nil
}

func (s *sSftp) Clear() error {
	if s.sftpClient != nil {
		_ = s.sftpClient.Close()
	}
	if s.sshClient != nil {
		return s.sshClient.Close()
	}
	return nil
}

func (s *sSftp) buildConfig() error {
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
	if s.Method == "" {
		s.Method = "password"
	}
	s.sshClientConfig, err = common.GetSSHClientConfig(
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
