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

	k := FractalMin + tools.SHA256Int(t.Data.SoloRings)%(FractalMax-FractalMin+1)
	if len(t.Data.SoloRings) < k {
		return nil
	}

	selectedRings := tools.RandomSelect(t.Data.SoloRings, k)
	team = t.selectVerificationTeam(selectedRings)
	if team == nil {
		return nil
	}
	t.Data.SoloRings = t.Data.SoloRings[k:]

	for i := 0; i < len(selectedRings); i++ {
		ring := t.Data.Rings[selectedRings[i]]

		ring.Next = selectedRings[(i+1)%k]
		ring.Prev = selectedRings[(i-1+k)%k]

		t.Data.Rings[selectedRings[i]] = ring
	}
	return
}

func (t *Trader) selectVerificationTeam(rings []string) []string {
	k := VerificationMin + tools.SHA256Int(rings)%(VerificationMax-VerificationMin+1)
	if len(t.Data.Traders) < k {
		return nil
	}

	traders := make([]string, 0, len(t.Data.Traders))
	for trader := range t.Data.Traders {
		traders = append(traders, trader)
	}
	return tools.RandomSelect(traders, k)
}
