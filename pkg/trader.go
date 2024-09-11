package pkg

import (
	"strconv"

	"github.com/Arka-Lab/LoR/tools"
)

type Trader struct {
	ID      string
	Account float64
	Wallet  string
}

func (trader *Trader) CreateCoin(amount float64, coinType uint) *CoinTable {
	return &CoinTable{
		ID:       tools.SHA3_256(trader.ID + "-" + strconv.Itoa(int(coinType))),
		Amount:   amount,
		Status:   Run,
		Type:     coinType,
		BindedOn: trader.ID,
		Owner:    trader.ID,
	}
}
