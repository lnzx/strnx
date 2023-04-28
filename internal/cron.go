package internal

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/jackc/pgx/v5"
	"github.com/lnzx/strnx/tools"
	"log"
	"sync"
	"time"
)

func StartAsync() {
	s := gocron.NewScheduler(time.UTC)
	s.SingletonModeAll()

	_, err := s.Every(10).Minutes().Do(dailyEarningsJob)
	if err != nil {
		log.Println(err)
	}
	_, err = s.Every(27).Minutes().Do(monthlyEarningsJob)
	if err != nil {
		log.Println(err)
	}

	s.StartAsync()
}

func dailyEarningsJob() {
	now := time.Now().UTC()
	log.Println("cron daily earnings started...")
	wallets, err := SelectWallets()
	if err != nil {
		log.Println(err)
		return
	}
	n := len(wallets)
	if n == 0 {
		log.Println("cron daily earnings wallets 0 skip")
		return
	}
	before := tools.GetBeforeDay(now)

	ch := make(chan *WalletResult, n)
	var wg sync.WaitGroup
	wg.Add(n)

	for _, wallet := range wallets {
		go func(addr string) {
			defer wg.Done()
			wr, e := FetchWalletEarnings(addr, before, now)
			if e != nil {
				log.Println(e)
				return
			}
			ch <- wr
			log.Printf("job daily ok %s %f %s\n", wr.Address, wr.Balance, time.Now().UTC().Sub(now).String())
		}(wallet.Address)
	}

	wg.Wait()
	close(ch)

	if len(ch) == 0 {
		log.Println("cron daily chan 0 skip update")
		return
	}

	batch := &pgx.Batch{}
	for wr := range ch {
		batch.Queue(UPDATE_WALLET_DAILY, wr.Balance, wr.NodeCount, wr.Address)
	}
	br := pool.SendBatch(context.Background(), batch)
	if err = br.Close(); err != nil {
		log.Println(err)
		return
	}
	log.Printf("cron daily earnings end %s\n", time.Now().UTC().Sub(now).String())
}

func monthlyEarningsJob() {
	now := time.Now().UTC()
	log.Println("cron monthly earnings started...")
	wallets, err := SelectWallets()
	if err != nil {
		log.Println(err)
		return
	}
	n := len(wallets)
	if n == 0 {
		log.Println("cron monthly earnings wallets 0 skip")
		return
	}
	start, end := tools.GetMonthRange(now)

	ch := make(chan *WalletResult, n)
	var wg sync.WaitGroup
	wg.Add(n)

	for _, wallet := range wallets {
		go func(addr string) {
			defer wg.Done()
			wr, e := FetchWalletEarnings(addr, start, end)
			if e != nil {
				log.Println(err)
				return
			}
			ch <- wr
			log.Printf("job monthly ok %s %f %s\n", addr, wr.Balance, time.Now().UTC().Sub(now).String())
		}(wallet.Address)
	}

	wg.Wait()
	close(ch)

	if len(ch) == 0 {
		log.Println("cron monthly chan 0 skip update")
		return
	}

	batch := &pgx.Batch{}
	for wr := range ch {
		batch.Queue(UPDATE_WALLET_BALANCE, wr.Balance, wr.Address)

		for _, earn := range wr.Earnings {
			amount := earn.FilAmount
			if amount == 0 {
				continue
			}
			date := tools.GetDate(earn.Timestamp)
			batch.Queue(UPSERT_DAILY_EARN, amount, wr.Address, date, amount)
		}
	}

	br := pool.SendBatch(context.Background(), batch)
	if err = br.Close(); err != nil {
		log.Println(err)
		return
	}
	log.Printf("cron monthly earnings end %s\n", time.Now().UTC().Sub(now).String())
}
