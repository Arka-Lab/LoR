package pkg

import (
	"github.com/Arka-Lab/LoR/tools"
)

const (
	FractalMin      = 50
	FractalMax      = 200
	VerificationMin = 20
	VerificationMax = 50
)

func (t *Trader) checkForFractalRings() (team []string) {
	if len(t.Data.SoloRings) < FractalMin {
		return nil
	}

	// TODO use a better random number generator
	k := FractalMin + tools.SHA256Int(t.Data.SoloRings)%(FractalMax-FractalMin+1)
	if len(t.Data.SoloRings) < k {
		return nil
	}

	// TODO use a better random number generator
	selectedRings, rest := tools.RandomSet(t.Data.SoloRings, k)
	team = t.selectVerificationTeam(selectedRings)
	if team == nil {
		return nil
	}

	for i := 0; i < len(selectedRings); i++ {
		ring := t.Data.Rings[selectedRings[i]]

		ring.Next = selectedRings[(i+1)%k]
		ring.Prev = selectedRings[(i-1+k)%k]

		t.Data.Rings[selectedRings[i]] = ring
	}
	t.Data.SoloRings = rest
	return
}

func (t *Trader) selectVerificationTeam(rings []string) (team []string) {
	k := VerificationMin + tools.SHA256Int(rings)%(VerificationMax-VerificationMin+1)
	if len(t.Data.Traders) < k {
		return nil
	}

	traders := make([]string, 0, len(t.Data.Traders))
	for trader := range t.Data.Traders {
		traders = append(traders, trader)
	}

	// TODO use a better random number generator
	team, _ = tools.RandomSet(traders, k)
	return
}
