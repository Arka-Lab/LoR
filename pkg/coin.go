package pkg

type Status int

const (
	Expired Status = iota - 2
	Blocked
	Run
	Paid
)

type CoinTable struct {
	ID       string  `json:"id"`
	Amount   float64 `json:"amount"`
	Status   Status  `json:"status"`
	Type     uint    `json:"type"`
	Next     string  `json:"next"`
	Prev     string  `json:"prev"`
	BindedOn string  `json:"binded_on"`
	Owner    string  `json:"owner"`
}
