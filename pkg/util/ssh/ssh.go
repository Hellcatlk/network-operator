package ssh

import (
	"net"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
)

func Run(address string, username string, password string, cmd *exec.Cmd) error {
	config := &ssh.ClientConfig{
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		User: username,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	config.SetDefaults()

	if !strings.Contains(address, ":") {
		address = address + ":22"
	}
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(cmd.String())
}
