package ssh

import (
	"bytes"
	"os"

	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"

	lua "github.com/yuin/gopher-lua"
)

func Execute(L *lua.LState) int {
	client := checkClient(L)
	args := L.CheckTable(2)
	command := args.RawGetString("command").String()

	session, err := makeSession(client)
	if err != nil {
		return sshError(L, err)
	}
	defer func() {
		_ = session.Close()
	}()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err = session.Run(command); err != nil {
		return sshError(L, err)
	}

	result := L.NewTable()
	L.SetField(result, `stdout`, lua.LString(stdout.String()))
	L.SetField(result, `stderr`, lua.LString(stderr.String()))
	L.Push(result)

	return 1
}

func makeSession(client *ssh.Client) (*ssh.Session, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func Copy(L *lua.LState) int {
	client := checkClient(L)
	args := L.CheckTable(2)
	source := args.RawGetString("source").String()
	dest := args.RawGetString("destination").String()

	session, err := makeSession(client)
	if err != nil {
		return sshError(L, err)
	}
	defer func() {
		_ = session.Close()
	}()

	if err = scp.CopyPath(source, dest, session); err != nil {
		return sshError(L, err)
	}

	L.Push(lua.LString("Copied!"))

	return 1
}

type sshConfig struct {
	user   string
	host   string
	port   string
	key    string
	config *ssh.ClientConfig
}

func Client(L *lua.LState) int {
	args := L.CheckTable(1)
	host := args.RawGetString("host").String()
	user := args.RawGetString("user").String()
	port := args.RawGetString("port").String()
	if port == "nil" {
		port = "22"
	}

	key := args.RawGetString("key").String()

	conn := sshConfig{host: host,
		user: user,
		port: port,
		key:  key,
	}

	client, err := conn.buildClient()

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	ud := L.NewUserData()
	ud.Value = client
	L.SetMetatable(ud, L.GetTypeMetatable(`ssh`))
	L.Push(ud)
	return 1
}

func (s *sshConfig) buildConfig() error {

	key, err := os.ReadFile(s.key)
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}

	s.config = &ssh.ClientConfig{
		User:            s.user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return nil
}

func (s *sshConfig) buildClient() (*ssh.Client, error) {
	err := s.buildConfig()
	if err != nil {
		return nil, err
	}

	host := s.host + ":" + s.port

	client, err := ssh.Dial("tcp", host, s.config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func checkClient(L *lua.LState) *ssh.Client {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*ssh.Client); ok {
		return v
	}
	L.ArgError(1, "This is not a ssh client")
	return nil
}

func sshError(L *lua.LState, e error) int {
	L.Push(lua.LNil)
	L.Push(lua.LString(e.Error()))
	return 2
}
