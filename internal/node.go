package internal

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lnzx/strnx/tools"
	"log"
	"os"
	"strings"
)

const (
	InsertNodeSql      = "INSERT INTO node(name,ip,bandwidth,traffic,price,renew) VALUES ($1,$2,$3,$4,$5,$6)"
	UpdateSysInfoSql   = "UPDATE node SET cpu=$1,ram=$2,disk=$3,traffic=$4,node_id=$5,type=$6 WHERE ip = $7"
	UpdateNodeStateSql = "UPDATE node SET state=$1 WHERE node_id = $2"
)

var HTTP_API_TOKEN = os.Getenv("HTTP_API_TOKEN")

func GetNodes(c *fiber.Ctx) error {
	nodes, err := SelectNodes()
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	return c.JSON(nodes)
}

func AddNodes(c *fiber.Ctx) error {
	node := new(Node)
	if err := c.BodyParser(node); err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}
	if err := tools.ValidateStruct(node); err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	_, err := pool.Exec(context.Background(), InsertNodeSql, node.Name, node.IP, node.Bandwidth, node.Traffic, node.Price, node.Renew)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	if SSH_USER != "" && SSH_PASS != "" {
		go UpdateNodeSysInfo(node.IP)
	}
	return nil
}

func UpdateNodeSysInfo(ip string) {
	sys, err := GetSysInfo(ip)
	if err != nil {
		log.Println("GetSysInfo error:", err)
		return
	}

	_, err = pool.Exec(context.Background(), UpdateSysInfoSql, sys.Cpu, sys.Ram, sys.Disk, sys.Traffic, sys.NodeId, sys.Version, ip)
	if err != nil {
		log.Println(err)
	}
}

func DeleteNodes(c *fiber.Ctx) error {
	ip := c.Query("ip")
	if ip == "" {
		return fiber.ErrBadRequest
	}
	ips := strings.Split(ip, ",")
	_, err := pool.Exec(context.Background(), "DELETE FROM node WHERE ip = ANY ($1)", ips)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	return nil
}

func Upgrade(c *fiber.Ctx) error {
	ip := c.FormValue("ip")
	if ip == "" {
		return fiber.ErrBadRequest
	}
	ips := strings.Split(ip, ",")
	cmd := fmt.Sprintf(UpgradeCmd, HTTP_API_TOKEN)

	for _, s := range ips {
		go func(host string) {
			_, e := Cmd(host, cmd, false)
			if e != nil {
				log.Println("Upgrade error:", e)
				return
			}
		}(s)
	}
	return nil
}
