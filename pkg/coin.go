package pkg

import (
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

func (t *Trader) SaveCoin(coin CoinTable) (*CooperationTable, []string, error) {
	if _, ok := t.Data.Coins[coin.ID]; ok {
		return nil, nil, ErrCoinAlreadyExist
	}
	if coin.Type >= uint(len(t.Data.RunCoins)) {
		return nil, nil, ErrInvalidCoinType
	}
	if _, ok := t.Data.Traders[coin.BindedOn]; !ok {
		return nil, nil, ErrTraderNotFound
	}
	if _, ok := t.Data.Traders[coin.Owner]; !ok {
		return nil, nil, ErrTraderNotFound
	}

	t.Data.Coins[coin.ID] = coin
	t.Data.RunCoins[coin.Type] = append(t.Data.RunCoins[coin.Type], coin.ID)

	ring := t.checkForCooperationRings()
	if ring != nil {
		if _, ok := t.Data.Rings[ring.ID]; ok {
			return nil, nil, ErrRingAlreadyExist
		}

		t.Data.Rings[ring.ID] = *ring
		t.Data.SoloRings = append(t.Data.SoloRings, ring.ID)
		return ring, t.checkForFractalRings(), nil
	}
	return nil, nil, nil
}
