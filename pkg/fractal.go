package pkg

import (
	"errors"
	"reflect"
	"slices"

	"github.com/Arka-Lab/LoR/tools"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/rand"
)

const (
	FractalMin      = 50
	FractalMax      = 200
	VerificationMin = 20
	VerificationMax = 50
)

type FractalRing struct {
	ID               string             `json:"id"`
	CooperationRings []CooperationTable `json:"cooperation_rings"`
	VerificationTeam []string           `json:"verification_team"`
	SoloRings        []string           `json:"-"`
}

func (t *Trader) checkForFractalRing() (fractal *FractalRing) {
	soloRings := make([]string, 0)
	for _, cooperation := range t.Data.Cooperations {
		if cooperation.Next == "" && cooperation.Prev == "" {
			soloRings = append(soloRings, cooperation.ID)
		}
	}

	selectedRing := selectFractalRing(soloRings, "")
	if selectedRing == nil {
		return nil
	}

	traders := maps.Keys(t.Data.Traders)
	team := selectVerificationTeam(traders, selectedRing, "")
	if team == nil {
		return nil
	}

	fractalID := tools.SHA256Str(selectedRing)
	selectedCooperations := make([]CooperationTable, len(selectedRing))
	for i, ringID := range selectedRing {
		cooperation := t.Data.Cooperations[ringID]
		cooperation.FractalID = fractalID
		cooperation.Next = selectedRing[(i+1)%len(selectedRing)]
		cooperation.Prev = selectedRing[(i-1+len(selectedRing))%len(selectedRing)]
		selectedCooperations[i] = cooperation
		t.Data.Cooperations[ringID] = cooperation
	}

	return &FractalRing{
		ID:               fractalID,
		CooperationRings: selectedCooperations,
		SoloRings:        soloRings,
		VerificationTeam: team,
	}
}

func (t *Trader) validateFractalRing(fractal *FractalRing) error {
	selectedRings := make([]string, 0, len(fractal.CooperationRings))
	for _, cooperation := range fractal.CooperationRings {
		if err := t.validateCooperationRing(cooperation); err != nil {
			return err
		}
		selectedRings = append(selectedRings, cooperation.ID)
	}
	traders := maps.Keys(t.Data.Traders)

	if fractal.ID != tools.SHA256Str(selectedRings) {
		return errors.New("invalid fractal ring id")
	} else if !reflect.DeepEqual(selectedRings, selectFractalRing(fractal.SoloRings, selectedRings[0])) {
		return errors.New("invalid selected cooperation ring")
	} else if !reflect.DeepEqual(fractal.VerificationTeam, selectVerificationTeam(traders, selectedRings, fractal.VerificationTeam[0])) {
		return errors.New("invalid verification team")
	}
	return nil
}

func selectFractalRing(soloRings []string, firstRing string) (result []string) {
	if len(soloRings) < FractalMin {
		return nil
	}

	copiedRings := make([]string, len(soloRings))
	copy(copiedRings, soloRings)
	slices.Sort(copiedRings)

	k := FractalMin + tools.SHA256Int(copiedRings)%(FractalMax-FractalMin+1)
	if len(soloRings) < k {
		return nil
	}
	result = make([]string, k)

	if firstRing != "" {
		result[0] = firstRing
		for i := 0; i < len(copiedRings); i++ {
			if copiedRings[i] == firstRing {
				copiedRings[i] = copiedRings[0]
				copiedRings = copiedRings[1:]
				break
			}
		}
	} else {
		index := rand.Intn(len(soloRings))
		result[0] = copiedRings[index]

		copiedRings[index] = copiedRings[0]
		copiedRings = copiedRings[1:]
	}

	rnd := make([]int, 0)
	for i := 1; i < k; i++ {
		if len(rnd) == 0 {
			rnd = tools.SHA256Arr(result)
		}
		index := rnd[0] % len(copiedRings)
		result[i], rnd = copiedRings[index], rnd[1:]

		copiedRings[index] = copiedRings[0]
		copiedRings = copiedRings[1:]
	}
	return
}
