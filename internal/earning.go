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
	"sort"
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
	start, end := tools.GetBeforeDayN(now, 30)
	metrics, err := fetchNodesEarning(start, end)
	if err != nil {
		log.Println(err)
		return
	}
	if len(metrics) == 0 {
		log.Println("cron FetchNodesEarningJob metrics 0 skip")
		return
	}

	tops := GetTopNodesEarning(metrics, 100)
	nodeMap, err := fetchNodesStatus()
	if err != nil {
		log.Println(err)
		return
	}

	batch := &pgx.Batch{}
	for _, node := range tops {
		if node.FilAmount == 0 {
			continue
		}
		status := nodeMap[node.NodeId]
		geo := status.Geoloc
		batch.Queue(UPSERT_EARN, node.NodeId, node.FilAmount, node.PayoutStatus,
			status.Speedtest.Isp, geo.Country, geo.City, geo.Region, status.Created)
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
			err = rsp.Body.Close()
			if err != nil {
				return nil, err
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

// GetTopNodesEarning 对filAmount从大到小排序,只取前n个
func GetTopNodesEarning(metrics []PerNodeMetrics, n int) []PerNodeMetrics {
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].FilAmount > metrics[j].FilAmount
	})
	return metrics[0:n]
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
			err = rsp.Body.Close()
			if err != nil {
				return nil, err
			}
		}
		return nil, err
	}

	bytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	var stats NodeStatus
	err = json.Unmarshal(bytes, &stats)
	if err != nil {
		return nil, err
	}

	return ConvertNodesToMap(stats.Nodes), nil
}

func ConvertNodesToMap(nodes []Status) map[string]Status {
	nodeMap := make(map[string]Status)
	for _, node := range nodes {
		nodeMap[node.Id] = node
	}
	return nodeMap
}
