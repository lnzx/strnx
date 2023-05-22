package internal

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/lnzx/strnx/tools"
	"log"
	"strings"
)

func AddPool(c *fiber.Ctx) error {
	traffic := c.FormValue("traffic")
	if traffic == "" {
		return errors.New("traffic cannot be empty")
	}
	ip := c.FormValue("ip")
	if ip == "" {
		return fiber.ErrBadRequest
	}
	ips := strings.Split(ip, ",")
	if len(ips) < 2 {
		return errors.New("a node cannot create a traffic pool")
	}
	err := addPool(tools.ToInt(traffic), ips)
	if err != nil {
		log.Println("addPool error:", err)
		return fiber.ErrInternalServerError
	}
	return nil
}

func addPool(traffic int, ips []string) (err error) {
	tx, err := pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if e := tx.Rollback(context.Background()); e != nil {
				log.Println("rollback:", e)
			}
		}
	}()

	var id int
	err = tx.QueryRow(context.Background(), "INSERT INTO pool(traffic) VALUES ($1) RETURNING id", traffic).Scan(&id)
	if err != nil {
		return err
	}
	log.Println("insert pool id:", id)
	batch := &pgx.Batch{}
	for _, s := range ips {
		batch.Queue("UPDATE node SET pool_id = $1 WHERE ip = $2", id, s)
	}
	br := tx.SendBatch(context.Background(), batch)
	if err = br.Close(); err != nil {
		return err
	}
	if err = tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
