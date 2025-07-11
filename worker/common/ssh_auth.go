package common

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type sshClientConfig struct {
	clientConf *ssh.ClientConfig
	method     string // 认证方式
	username   string // 用户名
	password   string // 密码
	privateKey string // 私钥/私钥路径
	keyPass    string // 私钥密码
}

func GetSSHClientConfig(method string, username string, password string, privateKey string, keyPass string) (*ssh.ClientConfig, error) {
	s := &sshClientConfig{
		method:     method,
		username:   username,
		password:   password,
		privateKey: privateKey,
		keyPass:    keyPass,
	}
	return s.getClientConf()
}

func (s *sshClientConfig) getClientConf() (*ssh.ClientConfig, error) {
	switch s.method {
	case "privateKey":
		if err := s.privateKeyAuth(); err != nil {
			return nil, err
		}
	case "privateKeyWithPassphrase":
		if err := s.privateKeyWithPassphraseAuth(); err != nil {
			return nil, err
		}
	case "sshAgent":
		if err := s.sshAgentAuth(); err != nil {
			return nil, err
		}
	case "password":
		if err := s.passwordKeyAuth(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported auth method: %s", s.method)
	}
	return s.clientConf, nil
}

func (s *sshClientConfig) privateKeyAuth() error {
	privateKey, err := os.ReadFile(s.privateKey)
	if err != nil {
		privateKey = []byte(s.privateKey)
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return err
	}
	s.clientConf = &ssh.ClientConfig{
		User: s.username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return nil
}

func (s *sshClientConfig) privateKeyWithPassphraseAuth() error {
	privateKey, err := os.ReadFile(s.privateKey)
	if err != nil {
		privateKey = []byte(s.privateKey)
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(s.keyPass))
	if err != nil {
		return err
	}
	s.clientConf = &ssh.ClientConfig{
		User: s.username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return nil
}
func (s *sshClientConfig) sshAgentAuth() error {
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return err
	}

	agentClient := agent.NewClient(conn)
	s.clientConf = &ssh.ClientConfig{
		User: s.username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agentClient.Signers),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return nil
}

func (s *sshClientConfig) passwordKeyAuth() error {
	s.clientConf = &ssh.ClientConfig{
		User: s.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return nil
}
