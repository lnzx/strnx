package internal

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/lnzx/strnx/tools"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	SELECT_BALANCE_NODES = "SELECT COALESCE(SUM(nodes[1]), 0) AS active,COALESCE(SUM(nodes[2]),0) AS inactive, COALESCE(CAST(SUM(balance) as DECIMAL(18,2)),0) AS balance FROM wallet"
	SELECT_DAILY_EARN    = "SELECT CAST(SUM(earnings) as DECIMAL(18,2)) as earnings FROM daily WHERE date >= $1 AND date <= $2 GROUP BY date ORDER BY date"
)

var X_CMC_PRO_API_KEY = os.Getenv("X-CMC_PRO_API_KEY")
var FILUSD float32 = 0

func Summary(c *fiber.Ctx) error {
	earnings, nodes, err := selectBalanceNodes()
	if err != nil {
		log.Println(err)
	}
	dailys, err := selectDailyEarns()
	if err != nil {
		log.Println(err)
	}
	groups, err := selectGroups()
	if err != nil {
		log.Println(err)
	}
	cost := SelectNodesCosts()
	return c.JSON(fiber.Map{
		"nodes":    nodes,
		"cost":     cost,
		"roi":      roi(earnings, cost),
		"earnings": earnings,
		"dailys":   dailys,
		"time":     time.Now().UTC().Format("2006-01-02 15:03:04"),
		"groups":   groups,
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

func selectGroups() (groups []Group, err error) {
	rows, err := pool.Query(context.Background(), "SELECT \"group\" AS \"name\",SUM(balance) AS balance FROM wallet GROUP BY \"group\" ORDER BY balance DESC")
	if err != nil {
		return
	}
	groups, err = pgx.CollectRows(rows, pgx.RowToStructByName[Group])
	return
}

func GetFilUsd() float32 {
	url := fmt.Sprintf("https://pro-api.coinmarketcap.com/v2/tools/price-conversion?amount=1&id=2280&convert=USD")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("GetFilUsd:", err)
		return 0
	}
	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", X_CMC_PRO_API_KEY)

	rsp, err := tools.Do(req)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer rsp.Body.Close()

	bytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		log.Println(err)
		return 0
	}

	var filUsd FilUsd
	err = json.Unmarshal(bytes, &filUsd)
	if err != nil {
		log.Println(err)
		return 0
	}
	if filUsd.Status.ErrorCode == 0 {
		FILUSD = filUsd.Data.Quote.USD.Price
		log.Println("GET FIL/USD:", FILUSD)
		return FILUSD
	}
	return 0
}

func roi(earnings, cost float32) int {
	earnings = earnings * FILUSD
	return int((earnings - cost) / cost * 100)
}
