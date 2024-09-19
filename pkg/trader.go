package pkg

import (
	"crypto/rsa"
	"strconv"

	"github.com/Arka-Lab/LoR/tools"
)

const (
	KeySize = 2048
)

type TraderData struct {
	Traders    map[string]Trader
	Coins      map[string]CoinTable
	Rings      map[string]CooperationTable
	PrivateKey *rsa.PrivateKey
	RunCoins   [][]string
	SoloRings  []string
}

type Trader struct {
	ID        string         `json:"id"`
	Account   float64        `json:"account"`
	Wallet    string         `json:"wallet"`
	PublicKey *rsa.PublicKey `json:"public_key"`
	Data      *TraderData    `json:"-"`
}

func CreateTrader(account float64, wallet string, coinTypeCount int) *Trader {
	privateKey, err := tools.GeneratePrivateKey(KeySize)
	if err != nil {
		return nil
	}

	return &Trader{
		ID:        tools.SHA256Str(wallet + "-" + strconv.Itoa(int(coinTypeCount))),
		Account:   account,
		Wallet:    wallet,
		PublicKey: &privateKey.PublicKey,
		Data: &TraderData{
			Traders:    make(map[string]Trader),
			Coins:      make(map[string]CoinTable),
			Rings:      make(map[string]CooperationTable),
			RunCoins:   make([][]string, coinTypeCount+1),
			SoloRings:  make([]string, 0),
			PrivateKey: privateKey,
		},
	}
}

func (t *Trader) SaveTrader(trader Trader) error {
	trader.Data = nil
	if _, ok := t.Data.Traders[trader.ID]; ok {
		return ErrTraderAlreadyExist
	}

	t.Data.Traders[trader.ID] = trader
	return nil
}
