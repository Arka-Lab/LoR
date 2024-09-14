package pkg

import (
	"strconv"

	"github.com/Arka-Lab/LoR/tools"
)

type TraderData struct {
	Traders     map[string]Trader
	Coins       map[string]CoinTable
	Rings       map[string]CooperationTable
	RunCoins    [][]string
	SoloRings   []string
	CoinCounter uint
}

type Trader struct {
	ID      string      `json:"id"`
	Account float64     `json:"account"`
	Wallet  string      `json:"wallet"`
	Data    *TraderData `json:"-"`
}

func CreateTrader(account float64, wallet string, coinTypeCount int) *Trader {
	return &Trader{
		ID:      tools.SHA256Str(wallet + "-" + strconv.Itoa(int(coinTypeCount))),
		Account: account,
		Wallet:  wallet,
		Data: &TraderData{
			CoinCounter: 0,
			Traders:     make(map[string]Trader),
			Coins:       make(map[string]CoinTable),
			Rings:       make(map[string]CooperationTable),
			RunCoins:    make([][]string, coinTypeCount+1),
			SoloRings:   make([]string, 0),
		},
	}
}

func (t *Trader) CreateCoin(amount float64, coinType uint) *CoinTable {
	t.Data.CoinCounter++

	return &CoinTable{
		ID:       tools.SHA256Str(t.ID + "-" + strconv.Itoa(int(coinType)) + "-" + strconv.Itoa(int(t.Data.CoinCounter))),
		Amount:   amount,
		Status:   Run,
		Type:     coinType,
		BindedOn: t.ID,
		Owner:    t.ID,
	}
}

func (t *Trader) SaveTrader(trader Trader) error {
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
	return ring, nil, nil
}
