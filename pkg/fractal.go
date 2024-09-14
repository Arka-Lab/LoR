package pkg

import "github.com/Arka-Lab/LoR/tools"

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
