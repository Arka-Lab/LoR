package pkg

import (
	"errors"
	"reflect"
	"slices"

	"github.com/Arka-Lab/LoR/tools"
	"golang.org/x/exp/maps"
)

const (
	FractalMin      = 50
	FractalMax      = 200
	VerificationMin = 20
	VerificationMax = 50
)

type FractalRing struct {
	ID               string              `json:"id"`
	CooperationRings []*CooperationTable `json:"cooperation_rings"`
	VerificationTeam []string            `json:"verification_team"`
	SoloRings        []string            `json:"-"`
}

func (t *Trader) checkForFractalRing() (fractal *FractalRing) {
	soloRings := maps.Keys(t.Data.Cooperations.SoloRings)
	slices.Sort(soloRings)

	selectedRing := selectFractalRing(soloRings)
	if selectedRing == nil {
		return nil
	}

	k := len(selectedRing)
	selectedCooperations := make([]*CooperationTable, k)
	for i, ringID := range selectedRing {
		coin := t.Data.Cooperations.RemoveSoloRing(ringID)
		selectedCooperations[i], coin.Next, coin.Prev = coin, selectedRing[(i+1)%k], selectedRing[(i-1+k)%k]
	}

	traders := maps.Keys(t.Data.Traders)
	slices.Sort(traders)

	team := selectVerificationTeam(traders, selectedRing)
	if team == nil {
		return nil
	}

	return &FractalRing{
		ID:               tools.SHA256Str(selectedRing),
		CooperationRings: selectedCooperations,
		SoloRings:        soloRings,
		VerificationTeam: team,
	}
}

func (t *Trader) validateFractalRing(fractal *FractalRing) error {
	runCoins := make([]map[string]*CoinTable, len(fractal.CooperationRings))
	for i := 0; i < len(runCoins); i++ {
		runCoins[i] = make(map[string]*CoinTable)
	}
	for _, coin := range t.Data.Coins.Coins {
		if coin.Status == Run {
			runCoins[coin.Type][coin.ID] = coin
		}
	}

	selectedRing := make([]string, len(fractal.CooperationRings))
	for i, cooperation := range fractal.CooperationRings {
		if err := t.validateCooperationRing(cooperation); err != nil {
			return err
		} else if cooperation.ID != tools.SHA256Str(selectCooperationRing(runCoins, cooperation.Investor)) {
			return errors.New("invalid cooperation ring id")
		}

		selectedRing[i] = cooperation.ID
		for i, coinID := range cooperation.Members {
			delete(runCoins[i], coinID)
		}
	}

	traders := maps.Keys(t.Data.Traders)
	slices.Sort(traders)

	if fractal.ID != tools.SHA256Str(selectedRing) {
		return errors.New("invalid fractal ring id")
	} else if !reflect.DeepEqual(selectedRing, selectFractalRing(fractal.SoloRings)) {
		return errors.New("invalid selected cooperation ring")
	} else if !reflect.DeepEqual(fractal.VerificationTeam, selectVerificationTeam(traders, selectedRing)) {
		return errors.New("invalid verification team")
	}
	return nil
}

func selectFractalRing(soloRings []string) (result []string) {
	if len(soloRings) < FractalMin {
		return nil
	}
	k := FractalMin + tools.SHA256Int(soloRings)%(FractalMax-FractalMin+1)
	if len(soloRings) < k {
		return nil
	}
	copiedRings := make([]string, len(soloRings))
	copy(copiedRings, soloRings)

	for _, index := range tools.RandomIndexes(len(copiedRings), k) {
		result = append(result, copiedRings[index])
		copiedRings[index] = copiedRings[0]
		copiedRings = copiedRings[1:]
	}
	return
}
