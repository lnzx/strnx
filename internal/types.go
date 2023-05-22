package internal

import "time"

type User struct {
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=6"`
}

type Wallet struct {
	Name    string  `json:"name" validate:"required,min=1"`
	Address string  `json:"address" validate:"required,len=41"`
	Balance float32 `json:"balance"`
	Daily   float32 `json:"daily"`
	Nodes   []int16 `json:"nodes"`
	Group   string  `json:"group"`
}

type Daily struct {
	Earnings float32 `json:"earnings"`
}

type Node struct {
	ID        int     `json:"id"`
	Name      string  `json:"name" validate:"required,min=1"`
	IP        string  `json:"ip" validate:"required,ipv4"`
	Bandwidth int     `json:"bandwidth"`
	Traffic   string  `json:"traffic"`
	Price     float32 `json:"price"`
	Renew     string  `json:"renew"`
	State     string  `json:"state"`
	Type      string  `json:"type"`
	NodeId    string  `json:"nodeId" db:"node_id"`
	CPU       int     `json:"cpu"`
	Ram       string  `json:"ram"`
	Disk      string  `json:"disk"`
	PoolId    int     `json:"poolId" db:"pool_id"`
}

type WalletResult struct {
	GlobalStats struct {
		TotalEarnings float32 `json:"totalEarnings"`
	} `json:"globalStats"`
	Earnings []Earning `json:"earnings"`
	Nodes    []struct {
		Count int16  `json:"count"`
		State string `json:"state"`
	} `json:"nodes"`
	PerNodeMetrics []PerNodeMetrics `json:"perNodeMetrics"`
	Address        string
}

func (w *WalletResult) NodeCounts() (active, others int16) {
	for _, node := range w.Nodes {
		if node.State == "active" {
			active += node.Count
		} else {
			others += node.Count
		}
	}
	return
}

type Earning struct {
	FilAmount float32 `json:"filAmount"`
	Timestamp string  `json:"timestamp"`
}

type NodeMetrics struct {
	PerNodeMetrics []PerNodeMetrics `json:"perNodeMetrics"`
}

type PerNodeMetrics struct {
	NodeId       string    `json:"nodeId"`
	FilAmount    float32   `json:"filAmount"`
	PayoutStatus string    `json:"payoutStatus"`
	Max          string    `json:"max"`
	Isp          string    `json:"isp"`
	Country      string    `json:"country"`
	City         string    `json:"city"`
	Region       string    `json:"region"`
	Created      time.Time `json:"created"`
}

type NodeStatsResult struct {
	Nodes []Status `json:"nodes"`
}

type Status struct {
	Id     string `json:"id"`
	State  string `json:"state"`
	Geoloc struct {
		Country string `json:"country"`
		City    string `json:"city"`
		Region  string `json:"region"`
	} `json:"geoloc"`
	Speedtest struct {
		Isp string `json:"isp"`
	} `json:"speedtest"`
	Created time.Time `json:"createdAt"`
}

type SysInfo struct {
	Disk    string
	Cpu     int
	Ram     string
	Traffic string
	NodeId  string
	Version string
}

type Group struct {
	Name    string  `json:"name"`
	Balance float32 `json:"balance"`
}
