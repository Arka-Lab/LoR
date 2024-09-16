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

func (t *Trader) SaveTrader(trader Trader) error {
	trader.Data = nil
	if _, ok := t.Data.Traders[trader.ID]; ok {
		return ErrTraderAlreadyExist
	}

	t.Data.Traders[trader.ID] = trader
	return nil
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
