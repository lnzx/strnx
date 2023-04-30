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
	SELECT_BALANCE_NODES = "SELECT COALESCE(SUM(nodes[1]), 0) AS active,COALESCE(SUM(nodes[2]),0) AS inactive, COALESCE(CAST(SUM(balance) as DECIMAL(18,2)),0) AS balance FROM wallet"
	SELECT_DAILY_EARN    = "SELECT date,CAST(SUM(earnings) as DECIMAL(18,2)) as earnings FROM daily WHERE date >= $1 AND date <= $2 GROUP BY date ORDER BY date"
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
		"earnings": earnings,
		"dailys":   dailys,
		"time":     time.Now().UTC().Format("2006-01-02 15:03:04"),
	})
}

func selectBalanceNodes() (earnings float32, nodes []int16, err error) {
	var active int16
	var inactive int16
	err = pool.QueryRow(context.Background(), SELECT_BALANCE_NODES).Scan(&active, &inactive, &earnings)
	if err != nil {
		return
	}
	nodes = append(nodes, active, inactive)
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
