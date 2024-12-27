package gossh

import (
	"fmt"
	"github.com/892294101/goscp/dec"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"log"
	"strings"
)

type SSHClient struct {
	config *dec.UploadStruct
	client *goph.Client
}

func NewSSHClient(config *dec.UploadStruct) (*SSHClient, error) {
	s := &SSHClient{
		config: config,
	}
	if err := s.connectSSH(config); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SSHClient) connectSSH(us *dec.UploadStruct) error {

	var c goph.Config
	c.User = us.HostUser
	c.Port = us.Port
	c.Addr = us.Host
	c.Callback = ssh.InsecureIgnoreHostKey()
	c.Auth = goph.Password(us.HostPass)
	client, err := goph.NewConn(&c)
	if err != nil {
		return err
	}

	s.client = client
	return nil
}

func (s *SSHClient) CreateDir(f string) error {
	if !strings.EqualFold(strings.ToLower(f), "enable") {
		return nil
	}
	if _, err := s.executeCommand(fmt.Sprintf("mkdir -p %s", strings.TrimSpace(s.config.TargetLocation))); err != nil {
		return err
	}
	if _, err := s.executeCommand(fmt.Sprintf("chmod -R 777 %s", strings.TrimSpace(s.config.TargetLocation))); err != nil {
		return err
	}
	return nil
}

func (s *SSHClient) executeCommand(c string) (string, error) {
	log.Printf("execute command: %s\n", c)
	out, err := s.client.Run(c)
	if err != nil {
		return "", fmt.Errorf("failed to execute %s command: %w", c, err)
	}
	return string(out), nil
}

func (s *SSHClient) Close() {
	s.client.Close()
}
