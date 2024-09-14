package pkg

import "github.com/Arka-Lab/LoR/tools"

const (
	RoundsCount     = 10
	FractalMin      = 500
	FractalMax      = 2000
	VerificationMin = 20
	VerificationMax = 50
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
