package internal

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5"
	"github.com/lnzx/strnx/tools"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	NodeStatusUrl   = "https://orchestrator.strn.pl/stats?sortColumn=id"
	NodesEarningUrl = "https://uc2x7t32m6qmbscsljxoauwoae0yeipw.lambda-url.us-west-2.on.aws/?filAddress=all&startDate=%d&endDate=%d&step=day"
	UPSERT_EARN     = "INSERT INTO earn(node_id,earning,status,isp,country,city,region,created) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT (node_id) DO UPDATE SET earning = $2"
)

// FetchNodesEarningJob Query the total earnings of all nodes in the last 30 days
func FetchNodesEarningJob() {
	now := time.Now().UTC()
	start, end := tools.GetMonthRange(now)
	metrics, err := fetchNodesEarning(start, end)
	if err != nil {
		log.Println(err)
		return
	}
	if len(metrics) == 0 {
		log.Println("cron FetchNodesEarningJob metrics 0 skip")
		return
	}

	batch := &pgx.Batch{}
	for _, node := range metrics {
		if node.FilAmount == 0 {
			continue
		}
		if status, ok := NodeStatusMap[node.NodeId]; ok {
			geo := status.Geoloc
			batch.Queue(UPSERT_EARN, node.NodeId, node.FilAmount, node.PayoutStatus,
				status.Speedtest.Isp, geo.Country, geo.City, geo.Region, status.Created)
		}
	}
	br := pool.SendBatch(context.Background(), batch)
	if err = br.Close(); err != nil {
		log.Println(err)
		return
	}
	log.Printf("cron FetchNodesEarningJob started %s\n", time.Now().UTC().Sub(now).String())
}

// fetchNodesEarning Get earnings of all nodes based on start and end time
func fetchNodesEarning(start, end time.Time) ([]PerNodeMetrics, error) {
	url := fmt.Sprintf(NodesEarningUrl, start.UnixMilli(), end.UnixMilli())
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

	bytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	var metrics NodeMetrics
	err = json.Unmarshal(bytes, &metrics)
	if err != nil {
		return nil, err
	}
	return metrics.PerNodeMetrics, nil
}

func fetchNodesStatus() (map[string]Status, error) {
	req, err := http.NewRequest("GET", NodeStatusUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	rsp, err := tools.Do(req)
	if err != nil {
		if rsp != nil {
			if e := rsp.Body.Close(); e != nil {
				log.Println(e)
			}
		}
		return nil, err
	}

	bytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	var r NodeStatsResult
	err = json.Unmarshal(bytes, &r)
	if err != nil {
		return nil, err
	}

	return ConvertNodesToMap(r.Nodes), nil
}

func ConvertNodesToMap(nodes []Status) map[string]Status {
	nodeMap := make(map[string]Status, len(nodes))
	for _, node := range nodes {
		nodeMap[node.Id] = node
	}
	return nodeMap
}
