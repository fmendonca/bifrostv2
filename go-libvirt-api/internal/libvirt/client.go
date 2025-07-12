package libvirt

import (
	"crypto/tls"
	"io/ioutil"
	"time"

	"golang.org/x/crypto/ssh"
	"libvirt.org/go/libvirt"
)

// TCP over TLS connection
func ConnectTCPTLS(address string) (*libvirt.Libvirt, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", address, conf)
	if err != nil {
		return nil, err
	}
	l := libvirt.New(conn)
	if err := l.Connect(); err != nil {
		return nil, err
	}
	return l, nil
}

// SSH tunnel connection
func ConnectSSH(address, user, keyPath string) (*libvirt.Libvirt, error) {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	client, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return nil, err
	}

	conn, err := client.Dial("unix", "/var/run/libvirt/libvirt-sock")
	if err != nil {
		return nil, err
	}

	l := libvirt.New(conn)
	if err := l.Connect(); err != nil {
		return nil, err
	}
	return l, nil
}
