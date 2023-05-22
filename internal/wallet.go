package internal

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/lnzx/strnx/tools"
	"io"
	"log"
	"strings"
	"time"
)

const WALLET_URL = "https://uc2x7t32m6qmbscsljxoauwoae0yeipw.lambda-url.us-west-2.on.aws/?filAddress=%s&startDate=%d&endDate=%d&step=day"

const (
	UPDATE_WALLET_DAILY   = "UPDATE wallet SET daily = $1, nodes = $2 WHERE address = $3"
	UPDATE_WALLET_BALANCE = "UPDATE wallet SET balance = $1 WHERE address = $2"
	UPSERT_DAILY_EARN     = "INSERT INTO daily(earnings,address,date) VALUES ($1, $2, $3) ON CONFLICT (address,date) DO UPDATE SET earnings = $1"
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

func FetchWalletEarnings(address string, start, end time.Time) (*WalletResult, error) {
	url := fmt.Sprintf(WALLET_URL, address, start.UnixMilli(), end.UnixMilli())
	rsp, err := tools.Get(url)
	if err != nil {
		if rsp != nil {
			e := rsp.Body.Close()
			if e != nil {
				log.Println(e)
			}
		}
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		e := Body.Close()
		if e != nil {
			log.Println(e)
		}
	}(rsp.Body)

	bytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	wr := &WalletResult{
		Address: address,
	}
	err = json.Unmarshal(bytes, wr)
	return wr, err
}
