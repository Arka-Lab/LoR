package pkg

import (
	"errors"
	"slices"

	"github.com/Arka-Lab/LoR/tools"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/rand"
)

const (
	RoundsCount = 10
	RoundLength = 100
)

type CooperationTable struct {
	ID       string   `json:"id"`
	Weight   float64  `json:"weight"`
	Next     string   `json:"next"`
	Prev     string   `json:"prev"`
	Investor string   `json:"investor"`
	Rounds   uint     `json:"rounds"`
	Members  []string `json:"-"`
}

func (t *Trader) checkForCooperationRing() *CooperationTable {
	for _, coins := range t.Data.Coins.RunCoins {
		if len(coins) == 0 {
			return nil
		}
	}

	selectedCoins := selectCooperationRing(t.Data.Coins.RunCoins, "")
	for i, coinID := range selectedCoins {
		coin := t.Data.Coins.RemoveRunCoin(coinID)
		coin.Next = selectedCoins[(i+1)%len(selectedCoins)]
		coin.Prev = selectedCoins[(i-1+len(selectedCoins))%len(selectedCoins)]
	}

	return &CooperationTable{
		ID:       tools.SHA256Str(selectedCoins),
		Weight:   t.calculateWeight(selectedCoins),
		Investor: selectedCoins[0],
		Rounds:   RoundsCount,
		Members:  selectedCoins,
	}
}

func (t *Trader) calculateWeight(ring []string) (weight float64) {
	for _, coinID := range ring {
		weight += t.Data.Coins.GetCoin(coinID).Amount
	}
	return
}

func (t *Trader) validateCooperationRing(cooperation *CooperationTable) error {
	weight := 0.0
	for i, coinID := range cooperation.Members {
		coin := t.Data.Coins.GetCoin(coinID)
		if coin == nil {
			return errors.New("coin not found")
		} else if coin.Status != Run {
			return errors.New("invalid coin status")
		} else if coin.Type != uint(i) {
			return errors.New("invalid coin type")
		} else if coin.Owner != coin.BindedOn {
			return errors.New("invalid coin owner/binded on")
		}
		weight += coin.Amount
	}

	if cooperation.ID != tools.SHA256Str(cooperation.Members) {
		return errors.New("invalid cooperation ring id")
	} else if cooperation.Weight != weight {
		return errors.New("invalid cooperation ring weight")
	} else if cooperation.Investor != cooperation.Members[0] {
		return errors.New("invalid cooperation ring investor")
	} else if cooperation.Rounds != RoundsCount {
		return errors.New("invalid cooperation ring rounds")
	}
	return nil
}

type CooperationSet struct {
	Cooperations map[string]*CooperationTable
	SoloRings    map[string]*CooperationTable
}

func NewCooperationSet() *CooperationSet {
	return &CooperationSet{
		Cooperations: make(map[string]*CooperationTable),
		SoloRings:    make(map[string]*CooperationTable),
	}
}

func (cs *CooperationSet) AddCooperationRing(cooperation *CooperationTable) error {
	if _, ok := cs.Cooperations[cooperation.ID]; ok {
		return errors.New("cooperation ring already exist")
	}

	cs.Cooperations[cooperation.ID] = cooperation
	cs.SoloRings[cooperation.ID] = cooperation
	return nil
}

func (cs *CooperationSet) GetCooperationRing(id string) *CooperationTable {
	if cooperation, ok := cs.Cooperations[id]; ok {
		return cooperation
	}
	return nil
}

func (cs *CooperationSet) RemoveSoloRing(ringID string) *CooperationTable {
	if cooperation := cs.GetCooperationRing(ringID); cooperation != nil {
		delete(cs.SoloRings, ringID)
		return cooperation
	}
	return nil
}

func selectCooperationRing(runCoins []map[string]*CoinTable, investor string) []string {
	rnd := make([]int, 0)
	selectedRing := make([]string, len(runCoins))
	if investor == "" {
		selectedRing[0] = maps.Keys(runCoins[0])[rand.Intn(len(runCoins[0]))]
	} else {
		selectedRing[0] = investor
	}
	for i := 1; i < len(runCoins); i++ {
		if len(rnd) == 0 {
			rnd = tools.SHA256Arr(selectedRing)
		}
		coins := maps.Keys(runCoins[i])
		slices.Sort(coins)

		rnd, selectedRing[i] = rnd[1:], coins[rnd[0]%len(runCoins[i])]
	}
	return selectedRing
}
