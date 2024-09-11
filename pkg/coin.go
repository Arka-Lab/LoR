package pkg

type Status int

const (
	Expired Status = iota - 2
	Blocked
	Run
	Paid
)

type CoinTable struct {
	ID       string
	Amount   float64
	Status   Status
	Type     uint
	Next     string
	Prev     string
	BindedOn string
	Owner    string
}
