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
	// 获取FIL/USD价格
	if X_CMC_PRO_API_KEY != "" {
		if _, err := s.Every(1).Days().Do(GetFilUsd); err != nil {
			log.Println(err)
		}
	}

	// 如果配置了发短信,才启动检查版本定时任务
	if smsApiKey != "" && mobile != "" {
		_, err := s.Every(5).Minutes().Do(CheckVersionJob)
		if err != nil {
			log.Println(err)
		}
	}

	// 如果配置了ssh用户名密码,才启动检查节点服务器信息定时任务
	if SSH_USER != "" && SSH_PASS != "" {
		_, err := s.Every(11).Minutes().Do(func() {
			time.Sleep(20 * time.Second)
			updateNodeSysInfoJob()
		})
		if err != nil {
			log.Println("node stats job error:", err)
		}
	}

	// 过去1天的收益
	if _, err := s.Every(13).Minutes().Do(dailyEarningsJob); err != nil {
		log.Println(err)
	}

	// 本月的收益
	if _, err := s.Every(30).Minutes().Do(func() {
		time.Sleep(1 * time.Minute)
		monthlyEarningsJob()
	}); err != nil {
		log.Println(err)
	}

	s.StartAsync()
}

func updateNodeSysInfoJob() {
	now := time.Now().UTC()
	nodes, err := SelectNodes()
	if err != nil {
		log.Println("cron updateNodeSysInfoJob:", err)
		return
	}
	if len(nodes) == 0 {
		log.Println("cron updateNodeSysInfoJob 0 skip")
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(nodes))
	batch := &pgx.Batch{}
	for _, node := range nodes {
		go func(ip string) {
			defer wg.Done()
			sys, e := GetSysInfo(ip)
			if e != nil {
				log.Println("job GetSysInfo error:", e)
				return
			}
			batch.Queue(UpdateSysInfoSql, sys.Cpu, sys.Ram, sys.Disk, sys.Traffic, sys.NodeId, sys.Version, ip)
		}(node.IP)
	}
	wg.Wait()

	if batch.Len() == 0 {
		log.Println("cron updateNodeSysInfoJob batch 0 skip")
		return
	}
	br := pool.SendBatch(context.Background(), batch)
	if err = br.Close(); err != nil {
		log.Println("cron updateNodeSysInfoJob", err)
		return
	}
	log.Printf("cron updateNodeInfoJob started %d %s\n", len(nodes), time.Now().UTC().Sub(now).String())
}

func dailyEarningsJob() {
	earningsJob("daily", tools.GetBeforeDay, func(wr *WalletResult, batch *pgx.Batch) {
		active, others := wr.NodeCounts()
		batch.Queue(UPDATE_WALLET_DAILY, wr.GlobalStats.TotalEarnings, []int16{active, others}, wr.Address)

		for _, metric := range wr.PerNodeMetrics {
			batch.Queue(UpdateNodeStateSql, metric.Max, metric.NodeId)
		}
	})
}

func monthlyEarningsJob() {
	earningsJob("monthly", tools.GetMonthRange, func(wr *WalletResult, batch *pgx.Batch) {
		batch.Queue(UPDATE_WALLET_BALANCE, wr.GlobalStats.TotalEarnings, wr.Address)

		for _, earn := range wr.Earnings {
			amount := earn.FilAmount
			if amount == 0 {
				continue
			}
			date := tools.GetDateString(earn.Timestamp)
			batch.Queue(UPSERT_DAILY_EARN, amount, wr.Address, date)
		}
	})
}

func earningsJob(name string, timeRangeFunc func(time.Time) (time.Time, time.Time), populateSqlFunc func(*WalletResult, *pgx.Batch)) {
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
	var wg sync.WaitGroup
	wg.Add(n)
	batch := &pgx.Batch{}
	for _, wallet := range wallets {
		go func(addr string) {
			defer wg.Done()
			wr, e := FetchWalletEarnings(addr, start, end)
			if e != nil {
				log.Println("job FetchWalletEarnings error:", e)
				return
			}
			populateSqlFunc(wr, batch)
			log.Printf("%s %s %f", name, addr, wr.GlobalStats.TotalEarnings)
		}(wallet.Address)
	}
	wg.Wait()

	if batch.Len() == 0 {
		log.Println("cron FetchWalletEarnings batch 0 skip")
		return
	}

	br := pool.SendBatch(context.Background(), batch)
	if err = br.Close(); err != nil {
		log.Println(err)
		return
	}
	log.Printf("cron %s earnings started %s\n", name, time.Now().UTC().Sub(now).String())
}
