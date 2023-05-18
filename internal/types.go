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
	CPU       int     `json:"cpu"`
	Ram       string  `json:"ram"`
	Disk      string  `json:"disk"`
}

type WalletResult struct {
	Earnings []Earning `json:"earnings"`
	Nodes    []struct {
		Count int16  `json:"count"`
		State string `json:"state"`
	} `json:"nodes"`
	Address string
}

func (w *WalletResult) Balance() (balance float32) {
	for _, earn := range w.Earnings {
		amount := earn.FilAmount
		if amount > 0 {
			balance += earn.FilAmount
		}
	}
	return
}

func (w *WalletResult) NodeCounts() (active, inactive int16) {
	for _, node := range w.Nodes {
		if node.State == "active" {
			active += node.Count
		} else {
			inactive += node.Count
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
	Isp          string    `json:"isp"`
	Country      string    `json:"country"`
	City         string    `json:"city"`
	Region       string    `json:"region"`
	Created      time.Time `json:"created"`
}

type NodeStatus struct {
	Nodes []Status `json:"nodes"`
}

type Status struct {
	Id     string `json:"id"`
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
	Version string
}
