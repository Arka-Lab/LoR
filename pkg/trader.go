package pkg

import (
	"strconv"

	"github.com/Arka-Lab/LoR/tools"
)

type Data struct {
	Traders map[string]*Trader
	Coins   map[string]*CoinTable
}

type Trader struct {
	ID      string  `json:"id"`
	Account float64 `json:"account"`
	Wallet  string  `json:"wallet"`
	Data    *Data   `json:"-"`
}

func (trader *Trader) CreateCoin(amount float64, coinType uint) *CoinTable {
	return &CoinTable{
		ID:       tools.SHA256(trader.ID + "-" + strconv.Itoa(int(coinType))),
		Amount:   amount,
		Status:   Run,
		Type:     coinType,
		BindedOn: trader.ID,
		Owner:    trader.ID,
	}
}
