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
	earningsJob("daily", tools.GetBeforeDay, func(wr *WalletResult, batch *pgx.Batch) {
		batch.Queue(UPDATE_WALLET_DAILY, wr.Balance, wr.NodeCount, wr.Address)
	})
}

func monthlyEarningsJob() {
	earningsJob("monthly", tools.GetMonthRange, func(wr *WalletResult, batch *pgx.Batch) {
		batch.Queue(UPDATE_WALLET_BALANCE, wr.Balance, wr.Address)

		for _, earn := range wr.Earnings {
			amount := earn.FilAmount
			if amount == 0 {
				continue
			}
			date := tools.GetDate(earn.Timestamp)
			batch.Queue(UPSERT_DAILY_EARN, amount, wr.Address, date, amount)
		}
	})
}

func earningsJob(name string, timeRangeFunc func(time.Time) (time.Time, time.Time), updateFunc func(*WalletResult, *pgx.Batch)) {
	now := time.Now().UTC()
	log.Printf("cron %s earnings started...\n", name)

	wallets, err := SelectWallets()
	if err != nil {
		log.Println(err)
		return
	}
	n := len(wallets)
	if n == 0 {
		log.Printf("cron %s earnings wallets 0 skip\n", name)
		return
	}

	start, end := timeRangeFunc(now)

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
		}(wallet.Address)
	}

	wg.Wait()
	close(ch)

	if len(ch) == 0 {
		log.Printf("cron %s earnings chan 0 skip\n", name)
		return
	}

	batch := &pgx.Batch{}
	for wr := range ch {
		updateFunc(wr, batch)
	}

	br := pool.SendBatch(context.Background(), batch)
	if err = br.Close(); err != nil {
		log.Println(err)
		return
	}
	log.Printf("cron %s earnings end %s\n", name, time.Now().UTC().Sub(now).String())
}
