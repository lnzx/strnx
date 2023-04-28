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
	Nodes   int     `json:"nodes"`
}

type WalletResult struct {
	Earnings  []Earning `json:"earnings"`
	Nodes     []Node    `json:"nodes"`
	Address   string
	Balance   float32
	NodeCount int
}

type Earning struct {
	FilAmount float32 `json:"filAmount"`
	Timestamp string  `json:"timestamp"`
}

type Node struct {
	Count int    `json:"count"`
	State string `json:"state"`
}

type Daily struct {
	Earnings float32 `json:"earnings"`
}
