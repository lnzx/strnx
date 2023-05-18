package internal

import (
	"bytes"
	"fmt"
	"github.com/lnzx/strnx/tools"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"strings"
	"time"
)

const (
	TIMEOUT    = 30 * time.Second
	PORT       = 22
	VersionCmd = "docker inspect saturn-node | grep -oP '(?<=\"VERSION=)[^_]+'"
	UpgradeCmd = "curl -H 'Authorization: Bearer %s' localhost:8080/v1/update"
	SysInfoCmd = "echo \"`df -h / | awk 'NR>1 {print $2, $3, $5}'` `nproc` `free -h | awk '/^Mem:/ {print $2}'`\""
)

var SSH_USER = os.Getenv("SSH_USER")
var SSH_PASS = os.Getenv("SSH_PASS")

func GetSysInfo(host string) (*SysInfo, error) {
	output, err := Cmd(host, SysInfoCmd, true)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// 99G 12G 13% 4 7.6Gi
	fields := strings.Fields(output)
	return &SysInfo{
		Disk: fields[1] + "/" + fields[0] + "(" + fields[2] + ")",
		Cpu:  tools.ToInt(fields[3]),
		Ram:  fields[4],
	}, nil
}

func GetVersion(host string) (string, error) {
	output, err := Cmd(host, VersionCmd, true)
	if err != nil {
		return "", err
	}
	return output, nil
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
		return "", err
	}
	defer func(client *ssh.Client) {
		e := client.Close()
		if e != nil {
			log.Println(e)
		}
	}(client)

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer func(session *ssh.Session) {
		e := session.Close()
		if e != nil {
			log.Println(e)
		}
	}(session)

	if result {
		var buf bytes.Buffer
		session.Stdout = &buf

		err = session.Run(cmd)
		if err != nil {
			return "", err
		}
		output := strings.TrimSpace(buf.String())
		return output, nil
	}

	err = session.Run(cmd)
	return "", err
}
