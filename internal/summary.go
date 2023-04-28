package internal

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/lnzx/strnx/tools"
	"log"
	"time"
)

const (
	SELECT_BALANCE_NODES = "SELECT SUM(nodes) AS nodes,CAST(SUM(balance) as DECIMAL(18,2)) AS balance FROM wallet"
	SELECT_DAILY_EARN    = "SELECT CAST(SUM(earnings) as DECIMAL(18,2)) as earnings FROM daily WHERE date >= $1 AND date <= $2 GROUP BY date ORDER BY date ASC"
)

func Summary(c *fiber.Ctx) error {
	earnings, nodes, err := selectBalanceNodes()
	if err != nil {
		log.Println(err)
	}
	dailys, err := selectDailyEarns()
	if err != nil {
		log.Println(err)
	}
	return c.JSON(fiber.Map{
		"nodes":    nodes,
		"inactive": 0,
		"earnings": earnings,
		"dailys":   dailys,
		"time":     time.Now().UTC().Format("2006-01-02 15:03:04"),
	})
}

func selectBalanceNodes() (earnings float32, nodes int, err error) {
	err = pool.QueryRow(context.Background(), SELECT_BALANCE_NODES).Scan(&nodes, &earnings)
	if err != nil {
		return
	}
	return
}

func selectDailyEarns() (dailys []float32, err error) {
	start, end := tools.GetMonthRangeDate(time.Now().UTC())
	rows, err := pool.Query(context.Background(), SELECT_DAILY_EARN, start, end)
	if err != nil {
		return
	}
	entries, err := pgx.CollectRows(rows, pgx.RowToStructByName[Daily])
	if err != nil {
		return
	}
	for _, entry := range entries {
		dailys = append(dailys, entry.Earnings)
	}
	return
}
