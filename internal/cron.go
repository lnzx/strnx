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

var NodeStatusMap map[string]Status

func StartAsync() {
	s := gocron.NewScheduler(time.UTC)
	// 如果配置了发短信,才启动检查版本定时任务
	if smsApiKey != "" && mobile != "" {
		_, err := s.Every(5).Minutes().Do(CheckVersionJob)
		if err != nil {
			log.Println(err)
		}
	}

	// 查询节点信息给其它任务使用[更新节点运行状态,保存所有节点收益]
	if _, err := s.Every(10).Minutes().Do(FetchNodeMapJob); err != nil {
		log.Println(err)
	}

	// 如果配置了ssh用户名密码,才启动检查节点服务器信息定时任务
	if SSH_USER != "" && SSH_PASS != "" {
		_, err := s.Every(13).Minutes().Do(func() {
			time.Sleep(20 * time.Second)
			updateNodeInfoJob()
		})
		if err != nil {
			log.Println("node stats job error:", err)
		}
	}

	// 过去1天的收益
	if _, err := s.Every(23).Minutes().Do(dailyEarningsJob); err != nil {
		log.Println(err)
	}

	// 本月的收益
	if _, err := s.Every(1).Hours().Do(func() {
		time.Sleep(1 * time.Minute)
		monthlyEarningsJob()
	}); err != nil {
		log.Println(err)
	}

	//if _, err := s.Every(1).Day().At("23:50;23:58").Do(FetchNodesEarningJob); err != nil {
	//	log.Println(err)
	//}

	s.StartAsync()
}

func FetchNodeMapJob() {
	status, err := fetchNodesStatus()
	if err != nil {
		log.Println("cron FetchNodeMap error:", err)
		return
	}
	NodeStatusMap = status
	log.Println("cron FetchNodeMap started:", len(status))
}

func updateNodeInfoJob() {
	nodes, err := SelectNodes()
	if err != nil {
		log.Println("cron updateNodeInfoJob error:", err)
		return
	}
	if len(nodes) == 0 {
		log.Println("cron updateNodeInfoJob 0, skip")
		return
	}

	for _, node := range nodes {
		go func(host, id string) {
			go UpdateNodeInfo(host)
		}(node.IP, node.NodeId)
	}
	log.Println("cron updateNodeInfoJob started:", len(nodes))
}

func dailyEarningsJob() {
	earningsJob("daily", tools.GetBeforeDay, func(wr *WalletResult, batch *pgx.Batch) {
		active, others := wr.NodeCounts()
		batch.Queue(UPDATE_WALLET_DAILY, wr.GlobalStats.TotalEarnings, []int16{active, others}, wr.Address)
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

	ch := make(chan *WalletResult, n)
	var wg sync.WaitGroup
	wg.Add(n)

	for _, wallet := range wallets {
		go func(addr string) {
			defer wg.Done()
			wr, e := FetchWalletEarnings(addr, start, end)
			if e != nil {
				log.Println("FetchWalletEarnings error:", e)
				return
			}
			log.Printf("%s %s %f", name, addr, wr.GlobalStats.TotalEarnings)
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
		populateSqlFunc(wr, batch)
	}

	br := pool.SendBatch(context.Background(), batch)
	if err = br.Close(); err != nil {
		log.Println(err)
		return
	}
	log.Printf("cron %s earnings started %s\n", name, time.Now().UTC().Sub(now).String())
}
