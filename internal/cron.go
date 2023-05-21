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
	// 如果配置了发短信,才启动检查版本定时任务
	if smsApiKey != "" && mobile != "" {
		_, err := s.Every(5).Minutes().Do(CheckVersionJob)
		if err != nil {
			log.Println(err)
		}
	}

	// 如果配置了ssh用户名密码,才启动检查节点服务器信息定时任务
	if SSH_USER != "" && SSH_PASS != "" {
		_, err := s.Every(13).Minutes().Do(nodeStatsJob)
		if err != nil {
			log.Println("node stats job error:", err)
		}
	}

	_, err := s.Every(10).Minutes().Do(dailyEarningsJob)
	if err != nil {
		log.Println(err)
	}
	_, err = s.Every(27).Minutes().Do(func() {
		time.Sleep(30 * time.Second)
		monthlyEarningsJob()
	})
	if err != nil {
		log.Println(err)
	}

	_, err = s.Every(1).Day().At("23:50;23:58").Do(FetchNodesEarningJob)
	if err != nil {
		log.Println(err)
	}

	s.StartAsync()
}

func nodeStatsJob() {
	nodes, err := SelectNodes()
	if err != nil {
		log.Println("node stats job error:", err)
		return
	}
	if len(nodes) == 0 {
		log.Println("node stats job 0, skip")
		return
	}

	for _, node := range nodes {
		go func(host string) {
			go UpdateSysInfo(host)
		}(node.IP)
	}
	log.Println("cron node stats job started:", len(nodes))
}

func dailyEarningsJob() {
	earningsJob("daily", tools.GetBeforeDay, func(wr *WalletResult, batch *pgx.Batch) {
		active, inactive := wr.NodeCounts()
		balance := wr.Balance()
		batch.Queue(UPDATE_WALLET_DAILY, balance, []int16{active, inactive}, wr.Address)
	})
}

func monthlyEarningsJob() {
	earningsJob("monthly", tools.GetMonthRange, func(wr *WalletResult, batch *pgx.Batch) {
		balance := wr.Balance()
		batch.Queue(UPDATE_WALLET_BALANCE, balance, wr.Address)

		for _, earn := range wr.Earnings {
			amount := earn.FilAmount
			if amount == 0 {
				continue
			}
			date := tools.GetDateString(earn.Timestamp)
			batch.Queue(UPSERT_DAILY_EARN, amount, wr.Address, date, amount)
		}
	})
}

func earningsJob(name string, timeRangeFunc func(time.Time) (time.Time, time.Time), updateFunc func(*WalletResult, *pgx.Batch)) {
	now := time.Now().UTC()
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
				log.Println(e)
				return
			} else {
				log.Printf("%s %s %f", name, addr, wr.Balance())
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
	log.Printf("cron %s earnings started %s\n", name, time.Now().UTC().Sub(now).String())
}
