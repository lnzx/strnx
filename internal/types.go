package internal

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

type WalletResult struct {
	Earnings []Earning `json:"earnings"`
	Nodes    []Node    `json:"nodes"`
	Address  string
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

type Node struct {
	Count int16  `json:"count"`
	State string `json:"state"`
}

type Daily struct {
	Earnings float32 `json:"earnings"`
}

type NodeMetrics struct {
	PerNodeMetrics []PerNodeMetrics `json:"perNodeMetrics"`
}

type PerNodeMetrics struct {
	NodeId       string  `json:"nodeId"`
	FilAmount    float32 `json:"filAmount"`
	PayoutStatus string  `json:"payoutStatus"`
}
