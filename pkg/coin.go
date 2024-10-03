package pkg

import (
	"errors"
	"fmt"
	"strconv"

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
	ID       string  `json:"id"`
	Amount   float64 `json:"amount"`
	Status   Status  `json:"status"`
	Type     uint    `json:"type"`
	Next     string  `json:"next"`
	Prev     string  `json:"prev"`
	BindedOn string  `json:"binded_on"`
	Owner    string  `json:"owner"`
}

func (t *Trader) CreateCoin(amount float64, coinType uint) *CoinTable {
	id, err := tools.SignWithPrivateKeyStr(t.ID+"-"+strconv.Itoa(int(coinType)), t.Data.PrivateKey)
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

func (t *Trader) SaveCoin(coin CoinTable) error {
	if coin.Status != Run {
		return errors.New("invalid coin status")
	} else if coin.Type >= t.Data.Coins.TypeCount {
		return errors.New("invalid coin type")
	} else if coin.BindedOn != coin.Owner {
		return errors.New("invalid coin binded on")
	} else if trader, ok := t.Data.Traders[coin.Owner]; !ok {
		return errors.New("trader not found")
	} else if err := tools.VerifyWithPublicKeyStr(coin.Owner+"-"+strconv.Itoa(int(coin.Type)), coin.ID, trader.PublicKey); err != nil {
		return errors.New("invalid coin id")
	}

	if err := t.Data.Coins.AddCoin(coin); err != nil {
		return err
	}
	cooperation := t.checkForCooperationRing()
	if cooperation != nil {
		if err := t.Data.Cooperations.AddCooperationRing(cooperation); err != nil {
			return err
		}
		fractal := t.checkForFractalRing()
		fmt.Println(fractal) // TODO: remove and complete the implementation
	}
	return nil
}

type CoinSet struct {
	TypeCount uint
	Coins     map[string]*CoinTable
	RunCoins  []map[string]*CoinTable
}

func NewCoinSet(typeCount uint) *CoinSet {
	set := &CoinSet{
		TypeCount: typeCount,
		Coins:     make(map[string]*CoinTable),
		RunCoins:  make([]map[string]*CoinTable, typeCount+1),
	}
	for i := 0; i <= int(typeCount); i++ {
		set.RunCoins[i] = make(map[string]*CoinTable)
	}
	return set
}

func (cs *CoinSet) GetCoin(id string) *CoinTable {
	if coin, ok := cs.Coins[id]; ok {
		return coin
	}
	return nil
}

func (cs *CoinSet) AddCoin(coin CoinTable) error {
	if _, ok := cs.Coins[coin.ID]; ok {
		return errors.New("coin already exist")
	}

	cs.RunCoins[coin.Type][coin.ID] = &coin
	cs.Coins[coin.ID] = &coin
	return nil
}

func (cs *CoinSet) RemoveRunCoin(coinID string) *CoinTable {
	if coin := cs.GetCoin(coinID); coin != nil {
		delete(cs.RunCoins[coin.Type], coinID)
		return coin
	}
	return nil
}
