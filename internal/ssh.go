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
	SysInfoCmd = `echo "$(nproc),$(free -h | awk '/^Mem:/ {print $2}'),$(df -h / | awk 'NR>1 {print $2, $3, $5}'),$(vnstat --oneline 2>/dev/null | awk -F ';' '{print $10}'),$(SATURN_HOME=$(sudo docker inspect saturn-node 2>/dev/null | grep 'SATURN_HOME' | sed 's/[" ,]//g' | awk -F '=' '{print $2}') && cat ${SATURN_HOME:-$HOME}/shared/nodeId.txt 2>/dev/null),$(sudo docker inspect saturn-node 2>/dev/null | grep -oP '(?<="VERSION=)[^_]+')"`
)

var SSH_USER = os.Getenv("SSH_USER")
var SSH_PASS = os.Getenv("SSH_PASS")

func GetSysInfo(host string) (*SysInfo, error) {
	output, err := Cmd(host, SysInfoCmd, true)
	if err != nil {
		log.Println("GetSysInfo error:", err)
		return nil, err
	}
	// 4,7.8Gi,158G 41G 27%,,eafb4495-470f-4a5b-8440-b44d043d85c3,887
	columns := strings.Split(output, ",")
	return &SysInfo{
		Cpu:     tools.ToInt(columns[0]),
		Ram:     columns[1],
		Disk:    columns[2],
		Traffic: tools.ParseTraffic(columns[3]),
		NodeId:  columns[4],
		Version: columns[5],
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
