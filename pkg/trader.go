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

func (t *Trader) checkForCooperationRings() *CooperationTable {
	for _, coins := range t.Data.RunCoins {
		if len(coins) == 0 {
			return nil
		}
	}

	members, selectedRing, weight := t.selectCooperationRing()

	return &CooperationTable{
		ID:          tools.SHA256Str(members),
		MemberCount: uint(len(selectedRing)),
		Weight:      weight,
		Investor:    selectedRing[0],
		Rounds:      RoundsCount,
		Members:     members,
	}
}

func (t *Trader) selectCooperationRing() (members []string, selectedRing []string, weight float64) {
	for index, coins := range t.Data.RunCoins {
		randomIndex := tools.SHA256Int(coins) % len(coins)
		selectedRing = append(selectedRing, coins[randomIndex])

		coins[randomIndex] = coins[len(coins)-1]
		t.Data.RunCoins[index] = coins[:len(coins)-1]
	}

	for i := 0; i < len(selectedRing); i++ {
		coin := t.Data.Coins[selectedRing[i]]

		weight += coin.Amount
		members = append(members, coin.ID)

		coin.Prev = selectedRing[(i-1+len(selectedRing))%len(selectedRing)]
		coin.Next = selectedRing[(i+1)%len(selectedRing)]
		coin.Status = Blocked

		t.Data.Coins[selectedRing[i]] = coin
	}

	return
}

func (t *Trader) checkForFractalRings() []string {
	if len(t.Data.SoloRings) < FractalMin {
		return nil
	}

	k := FractalMin + tools.SHA256Int(t.Data.SoloRings)%(FractalMax-FractalMin+1)
	if len(t.Data.SoloRings) < k {
		return nil
	}

	selectedRings, rest := tools.RandomSet(t.Data.SoloRings, k)
	for i := 0; i < len(selectedRings); i++ {
		ring := t.Data.Rings[selectedRings[i]]

		ring.Next = selectedRings[(i+1)%k]
		ring.Prev = selectedRings[(i-1+k)%k]

		t.Data.Rings[selectedRings[i]] = ring
	}
	t.Data.SoloRings = rest

	return t.selectVerificationTeam(selectedRings)
}

func (t *Trader) selectVerificationTeam(rings []string) (team []string) {
	k := VerificationMin + tools.SHA256Int(rings)%(VerificationMax-VerificationMin+1)
	selectedRings, _ := tools.RandomSet(rings, k)
	for _, ring := range selectedRings {
		selectedMember, _ := tools.RandomSet(t.Data.Rings[ring].Members, 1)
		team = append(team, selectedMember[0])
	}
	return
}
