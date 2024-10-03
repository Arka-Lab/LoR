package pkg

import (
	"crypto/rsa"
	"errors"
	"strconv"
	"time"

	"github.com/Arka-Lab/LoR/tools"
)

const (
	KeySize = 2048
)

type TraderData struct {
	Coins        *CoinSet
	Cooperations *CooperationSet
	Traders      map[string]Trader
	PrivateKey   *rsa.PrivateKey
	Ticker       *time.Ticker
}

type Trader struct {
	ID        string         `json:"id"`
	Account   float64        `json:"account"`
	Wallet    string         `json:"wallet"`
	PublicKey *rsa.PublicKey `json:"public_key"`
	Data      *TraderData    `json:"-"`
}

func CreateTrader(account float64, wallet string, coinTypeCount uint) *Trader {
	privateKey, err := tools.GeneratePrivateKey(KeySize)
	if err != nil {
		return nil
	}

	time.Sleep(time.Until(time.Now().Truncate(time.Second).Add(time.Second)))
	ticker := time.NewTicker(RoundLength * time.Millisecond)

	return &Trader{
		ID:        tools.SHA256Str(wallet + "-" + strconv.Itoa(int(coinTypeCount))),
		Account:   account,
		Wallet:    wallet,
		PublicKey: &privateKey.PublicKey,
		Data: &TraderData{
			Coins:        NewCoinSet(coinTypeCount),
			Cooperations: NewCooperationSet(),
			Traders:      make(map[string]Trader),
			PrivateKey:   privateKey,
			Ticker:       ticker,
		},
	}
}

func (t *Trader) SaveTrader(trader Trader) error {
	trader.Data = nil
	if _, ok := t.Data.Traders[trader.ID]; ok {
		return errors.New("trader already exist")
	} else if trader.ID != tools.SHA256Str(trader.Wallet+"-"+strconv.Itoa(int(t.Data.Coins.TypeCount))) {
		return errors.New("invalid trader ID")
	}

	t.Data.Traders[trader.ID] = trader
	return nil
}

// submission:
// 1. confirmation of verification team (DONE)
// 2. confirmation of fractal ring (DONE)
// 3. lock mechanism

// round checking:
// 1. get acknowledge from all traders at the end of the round
// 2. if not, then the round will be expired

// payment rules:
// 1. run at each checkpoint
// 2. if the coin status is "RUN", then the payment will be done and status becomes "PAID"/"TERMINATED"
