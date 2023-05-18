package internal

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lnzx/strnx/tools"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

const (
	TIMEOUT    = 30 * time.Second
	PORT       = 22
	UpgradeCmd = "curl -H 'Authorization: Bearer %s' localhost:8080/v1/update"
	CpuCmd     = "nproc"
	RamCmd     = "free -h | awk '/^Mem:/ {print $2}'"
	DiskCmd    = "df -h / | awk 'NR>1 {print $2, $3, $5}'"
	VersionCmd = "docker inspect saturn-node | grep -oP '(?<=\"VERSION=)[^_]+'"
	SysInfoCmd = CpuCmd + " && " + RamCmd + " && " + DiskCmd + " && " + VersionCmd
)

var SSH_USER = os.Getenv("SSH_USER")
var SSH_PASS = os.Getenv("SSH_PASS")

func GetSysInfo(host string) (*SysInfo, error) {
	output, err := Cmd(host, SysInfoCmd, true)
	if err != nil {
		log.Println("GetSysInfo error:", err)
		return nil, err
	}
	// 1
	// 961Mi
	// 49G 4.1G 9%
	// Error: No such object: saturn-node
	lines := strings.Split(output, "\n")
	return &SysInfo{
		Cpu:     tools.ToInt(lines[0]),
		Ram:     lines[1],
		Disk:    tools.WrapDisk(lines[2]),
		Version: tools.IfThen(strings.HasPrefix(lines[3], "Error:"), "-", lines[3]),
	}, nil
}

func Cmd(host, cmd string, result bool) (string, error) {
	cfg := &ssh.ClientConfig{
		User:            SSH_USER,
		Auth:            []ssh.AuthMethod{ssh.Password(SSH_PASS)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         TIMEOUT,
	}
	addr := fmt.Sprintf("%s:%d", host, PORT)
	client, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return "", fmt.Errorf("unable to connect: %v", err)
	}
	defer func() {
		if err = client.Close(); err != nil {
			log.Println("unable to close client:", err)
		}
	}()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("unable to create session: %v", err)
	}
	defer func() {
		if err = session.Close(); err != nil && !errors.Is(err, io.EOF) {
			log.Println("unable to close session:", err)
		}
	}()

	if result {
		var stdout bytes.Buffer
		session.Stdout = &stdout

		err = session.Run(cmd)
		if err != nil {
			return "", fmt.Errorf("unable to run command %q: %v", cmd, err)
		}
		return strings.TrimSpace(stdout.String()), nil
	}

	err = session.Run(cmd)
	if err != nil {
		err = fmt.Errorf("unable to run command %q: %v", cmd, err)
	}
	return "", err
}
