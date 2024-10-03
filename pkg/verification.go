package pkg

import (
	"github.com/Arka-Lab/LoR/tools"
)

func (t *Trader) SubmitRing(ring *FractalRing) error {
	if err := t.validateFractalRing(ring); err != nil {
		return err
	}
	return nil
}

func selectVerificationTeam(traders []string, ring []string) (team []string) {
	k := VerificationMin + tools.SHA256Int(ring)%(VerificationMax-VerificationMin+1)
	if len(traders) < k {
		return nil
	}

	for _, index := range tools.RandomIndexes(len(traders), k) {
		team = append(team, traders[index])
		traders[index] = traders[0]
		traders = traders[1:]
	}
	return
}
