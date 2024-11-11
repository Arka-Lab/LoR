package pkg

import (
	"errors"
	"fmt"

	"github.com/Arka-Lab/LoR/tools"
)

type Status int

const (
	Expired Status = iota - 2
	Blocked
	Run
	Paid
)

type CoinTable struct {
	ID            string  `json:"id"`
	Amount        float64 `json:"amount"`
	Status        Status  `json:"status"`
	Type          uint    `json:"type"`
	Next          string  `json:"next"`
	Prev          string  `json:"prev"`
	BindedOn      string  `json:"binded_on"`
	Owner         string  `json:"owner"`
	CooperationID string  `json:"-"`
}

func (t *Trader) CreateCoin(amount float64, coinType uint) *CoinTable {
	id, err := tools.SignWithPrivateKeyStr(t.ID+"-"+fmt.Sprint(coinType), t.Data.PrivateKey)
	if err != nil {
		return nil
	}

	return &CoinTable{
		ID:       id,
		Amount:   amount,
		Status:   Run,
		Type:     coinType,
		BindedOn: t.ID,
		Owner:    t.ID,
	}
}

func (t *Trader) SaveCoin(coin CoinTable, tryToRing bool) (*FractalRing, error) {
	if coin.Status != Run {
		return nil, errors.New("invalid coin status")
	} else if coin.Type >= t.Data.CoinTypeCount {
		return nil, errors.New("invalid coin type")
	} else if coin.BindedOn != coin.Owner {
		return nil, errors.New("invalid coin binded on")
	} else if trader, ok := t.Data.Traders[coin.Owner]; !ok {
		return nil, errors.New("trader not found")
	} else if err := tools.VerifyWithPublicKeyStr(coin.Owner+"-"+fmt.Sprint(coin.Type), coin.ID, trader.PublicKey); err != nil {
		return nil, errors.New("invalid coin id")
	} else if coin.Next != "" || coin.Prev != "" {
		return nil, errors.New("coin is already in a ring")
	} else if _, ok := t.Data.Coins[coin.ID]; ok {
		return nil, errors.New("coin already exist")
	}

	t.Data.Coins[coin.ID] = coin
	if tryToRing {
		return t.CheckForRings(), nil
	}
	return nil, nil
}
