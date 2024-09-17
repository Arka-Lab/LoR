package pkg

import (
	"math/rand"

	"github.com/Arka-Lab/LoR/tools"
)

const (
	RoundsCount = 10
)

type CooperationTable struct {
	ID          string   `json:"id"`
	MemberCount uint     `json:"member_count"`
	Weight      float64  `json:"weight"`
	Next        string   `json:"next"`
	Prev        string   `json:"prev"`
	Investor    string   `json:"investor"`
	Rounds      uint     `json:"rounds"`
	Members     []string `json:"-"`
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
	selectedRing = make([]string, len(t.Data.RunCoins))
	startIndex := rand.Intn(len(selectedRing))

	randomIndex := rand.Intn(len(t.Data.RunCoins[startIndex]))
	selectedRing[startIndex] = t.Data.RunCoins[startIndex][randomIndex]
	t.removeRunCoin(startIndex, randomIndex)
	for i := startIndex + 1; i < len(t.Data.RunCoins); i++ {
		randomIndex = tools.SHA256Int(selectedRing) % len(t.Data.RunCoins[i])
		selectedRing[i] = t.Data.RunCoins[i][randomIndex]
		t.removeRunCoin(i, randomIndex)
	}
	for i := 0; i < startIndex; i++ {
		randomIndex = tools.SHA256Int(selectedRing) % len(t.Data.RunCoins[i])
		selectedRing[i] = t.Data.RunCoins[i][randomIndex]
		t.removeRunCoin(i, randomIndex)
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

func (t *Trader) removeRunCoin(coinType int, index int) {
	if index < len(t.Data.RunCoins[coinType]) {
		t.Data.RunCoins[coinType][index] = t.Data.RunCoins[coinType][0]
		t.Data.RunCoins[coinType] = t.Data.RunCoins[coinType][1:]
	}
}
