package pkg

import (
	"fmt"
	"strconv"

	"github.com/Arka-Lab/LoR/tools"
)

type Data struct {
	Traders     map[string]Trader
	Coins       map[string]CoinTable
	RunCoins    [][]string
	CoinCounter uint
}

type Trader struct {
	ID      string  `json:"id"`
	Account float64 `json:"account"`
	Wallet  string  `json:"wallet"`
	Data    *Data   `json:"-"`
}

func CreateTrader(account float64, wallet string, coinTypeCount int) *Trader {
	return &Trader{
		ID:      tools.SHA256str(wallet + "-" + strconv.Itoa(int(coinTypeCount))),
		Account: account,
		Wallet:  wallet,
		Data: &Data{
			CoinCounter: 0,
			Traders:     make(map[string]Trader),
			Coins:       make(map[string]CoinTable),
			RunCoins:    make([][]string, coinTypeCount+1),
		},
	}
}

func (t *Trader) CreateCoin(amount float64, coinType uint) *CoinTable {
	t.Data.CoinCounter++

	return &CoinTable{
		ID:       tools.SHA256str(t.ID + "-" + strconv.Itoa(int(coinType)) + "-" + strconv.Itoa(int(t.Data.CoinCounter))),
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

func (t *Trader) SaveCoin(coin CoinTable) (*CooperationTable, error) {
	if _, ok := t.Data.Coins[coin.ID]; ok {
		return nil, ErrCoinAlreadyExist
	}

	t.Data.Coins[coin.ID] = coin
	t.Data.RunCoins[coin.Type] = append(t.Data.RunCoins[coin.Type], coin.ID)
	return t.checkForCooperationRings(), nil
}

func (t *Trader) checkForCooperationRings() *CooperationTable {
	for _, coins := range t.Data.RunCoins {
		if len(coins) == 0 {
			return nil
		}
	}

	ringMask, selectedRing, weight := t.selectRing()

	return &CooperationTable{
		ID:          tools.SHA256str(ringMask),
		MemberCount: uint(len(selectedRing)),
		Weight:      weight,
		Investor:    selectedRing[0],
		Rounds:      RoundsCount,
	}
}

func (t *Trader) selectRing() (ringMask string, selectedRing []string, weight float64) {
	for index, coins := range t.Data.RunCoins {
		randomIndex := tools.SHA256int(fmt.Sprint(coins)) % len(coins)
		selectedRing = append(selectedRing, coins[randomIndex])

		coins[randomIndex] = coins[len(coins)-1]
		t.Data.RunCoins[index] = coins[:len(coins)-1]
	}

	for i := 0; i < len(selectedRing); i++ {
		coin := t.Data.Coins[selectedRing[i]]
		weight += coin.Amount
		if i > 0 {
			ringMask += "-"
		}
		ringMask += coin.ID

		coin.Prev = selectedRing[(i-1+len(selectedRing))%len(selectedRing)]
		coin.Next = selectedRing[(i+1)%len(selectedRing)]
		coin.Status = Blocked

		t.Data.Coins[selectedRing[i]] = coin
	}

	return
}
