package internal

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/lnzx/strnx/tools"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const WALLET_URL = "https://uc2x7t32m6qmbscsljxoauwoae0yeipw.lambda-url.us-west-2.on.aws/?filAddress=%s&startDate=%d&endDate=%d&step=day"

const (
	UPDATE_WALLET_DAILY   = "UPDATE wallet SET daily = $1, nodes = $2 WHERE address = $3"
	UPDATE_WALLET_BALANCE = "UPDATE wallet SET balance = $1 WHERE address = $2"
	UPSERT_DAILY_EARN     = "INSERT INTO daily(earnings,address,date) VALUES ($1, $2, $3) ON CONFLICT (address,date) DO UPDATE SET earnings = $4"
)

func GetWallets(c *fiber.Ctx) error {
	wallets, err := SelectWallets()
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(wallets)
}

func AddWallet(c *fiber.Ctx) error {
	wallet := new(Wallet)
	if err := c.BodyParser(wallet); err != nil {
		return fiber.ErrBadRequest
	}
	if err := tools.ValidateStruct(wallet); err != nil {
		return fiber.ErrBadRequest
	}
	if err := InsertWallet(wallet); err != nil {
		return fiber.ErrInternalServerError
	}
	return nil
}

func DelWallets(c *fiber.Ctx) error {
	addr := c.Query("addrs")
	if addr == "" {
		return fiber.ErrBadRequest
	}
	addrs := strings.Split(addr, ",")
	if len(addrs) == 0 {
		return fiber.ErrBadRequest
	}
	err := Batch("DELETE FROM wallet WHERE address = $1", addrs)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	return nil
}

// FetchWalletEarnings Nodes return nil
func FetchWalletEarnings(address string, start, end time.Time) (*WalletResult, error) {
	url := fmt.Sprintf(WALLET_URL, address, start.UnixMilli(), end.UnixMilli())
	rsp, err := http.Get(url)
	if err != nil {
		if rsp != nil {
			rsp.Body.Close()
		}
		return nil, err
	}
	defer rsp.Body.Close()

	bytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	wr := &WalletResult{
		Address: address,
	}
	err = json.Unmarshal(bytes, wr)
	if err != nil {
		return nil, err
	}
	wr.Balance = calculateTotalFilAmount(wr.Earnings)
	wr.NodeCount = countTotalNodes(wr.Nodes)
	wr.Nodes = nil
	return wr, nil
}

func calculateTotalFilAmount(earnings []Earning) (total float32) {
	for _, earning := range earnings {
		total += earning.FilAmount
	}
	return total
}

func countTotalNodes(nodes []Node) (count int) {
	for _, node := range nodes {
		count += node.Count
	}
	return count
}
