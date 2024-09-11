package pkg_test

import (
	"testing"

	"github.com/Arka-Lab/LoR/pkg"
	"github.com/Arka-Lab/LoR/tools"
)

func TestCreateTrader(t *testing.T) {
	trader := pkg.CreateTrader(100, "test_wallet", 10)

	if trader.ID != tools.SHA256str("test_wallet-10") {
		t.Error("CreateTrader failed: ID")
	} else if trader.Account != 100 {
		t.Error("CreateTrader failed: Account")
	} else if trader.Wallet != "test_wallet" {
		t.Error("CreateTrader failed: Wallet")
	} else if len(trader.Data.Traders) != 0 {
		t.Error("CreateTrader failed: Traders")
	} else if len(trader.Data.Coins) != 0 {
		t.Error("CreateTrader failed: Coins")
	} else if len(trader.Data.RunCoins) != 11 {
		t.Error("CreateTrader failed: RunCoins")
	}
}

func TestCreateCoin(t *testing.T) {
	trader := pkg.CreateTrader(100, "test_wallet", 10)

	if coin := trader.CreateCoin(100, 1); coin.ID != tools.SHA256str(trader.ID+"-1-1") {
		t.Error("CreateCoin failed: ID")
	} else if coin.Amount != 100 {
		t.Error("CreateCoin failed: Amount")
	} else if coin.Status != pkg.Run {
		t.Error("CreateCoin failed: Status")
	} else if coin.Type != 1 {
		t.Error("CreateCoin failed: Type")
	} else if coin.BindedOn != trader.ID {
		t.Error("CreateCoin failed: BindedOn")
	} else if coin.Owner != trader.ID {
		t.Error("CreateCoin failed: Owner")
	}

	if coin := trader.CreateCoin(96.5, 1); coin.ID != tools.SHA256str(trader.ID+"-1-2") {
		t.Error("CreateCoin failed: ID")
	}
}

func TestSaveTrader(t *testing.T) {
	trader1 := pkg.CreateTrader(100, "test_wallet1", 10)
	trader2 := pkg.CreateTrader(100, "test_wallet2", 10)

	if err := trader1.SaveTrader(*trader2); err != nil {
		t.Fatal("SaveTrader failed:", err)
	}
	if err := trader1.SaveTrader(*trader2); err == nil {
		t.Fatal("SaveTrader failed: Fail to detect duplicate trader")
	}
}

func TestSaveCoin(t *testing.T) {
	trader1 := pkg.CreateTrader(100, "test_wallet1", 2)
	trader2 := pkg.CreateTrader(100, "test_wallet2", 2)

	coin1 := trader1.CreateCoin(10.5, 0)
	coin2 := trader2.CreateCoin(9.8, 1)
	coin3 := trader2.CreateCoin(7.3, 0)
	coin4 := trader2.CreateCoin(2.5, 1)
	coin5 := trader2.CreateCoin(3.4, 2)
	coin6 := trader1.CreateCoin(7.3, 2)
	coin7 := trader1.CreateCoin(11.7, 2)
	ringIDs, ringWeights, ringInvestors := []string{}, []float64{}, []string{}
	for _, c1 := range []*pkg.CoinTable{coin1, coin3} {
		for _, c2 := range []*pkg.CoinTable{coin2, coin4} {
			for _, c3 := range []*pkg.CoinTable{coin5, coin6} {
				ringIDs = append(ringIDs, tools.SHA256str(c1.ID+"-"+c2.ID+"-"+c3.ID))
				ringWeights = append(ringWeights, c1.Amount+c2.Amount+c3.Amount)
				ringInvestors = append(ringInvestors, c1.ID)
			}
		}
	}

	saveAndCheck := func(trader *pkg.Trader, coin *pkg.CoinTable, hasError, hasRing bool) (ring *pkg.CooperationTable) {
		ring, err := trader.SaveCoin(*coin)
		if err != nil && !hasError {
			t.Fatal("SaveCoin failed:", err)
		} else if err == nil && hasError {
			t.Fatal("SaveCoin failed: Error not detected")
		} else if ring != nil && !hasRing {
			t.Fatal("SaveCoin failed: Cooperation ring detected")
		} else if ring == nil && hasRing {
			t.Fatal("SaveCoin failed: Fail to detect cooperation ring")
		}

		return ring
	}
	findIndex := func(ringID string) int {
		for i, id := range ringIDs {
			if id == ringID {
				return i
			}
		}
		return -1
	}

	for _, trader := range []*pkg.Trader{trader1, trader2} {
		for _, coin := range []*pkg.CoinTable{coin1, coin2, coin3, coin4} {
			saveAndCheck(trader, coin, false, false)
		}

		ring1 := saveAndCheck(trader, coin5, false, true)
		ring2 := saveAndCheck(trader, coin6, false, true)
		if ring1.ID == ring2.ID {
			t.Fatal("SaveCoin failed: Same cooperation ring ID")
		}

		index1 := findIndex(ring1.ID)
		index2 := findIndex(ring2.ID)
		if index1 == -1 || index2 == -1 {
			t.Fatal("SaveCoin failed: Cooperation ring not found")
		} else if ring1.MemberCount != 3 || ring2.MemberCount != 3 {
			t.Fatal("SaveCoin failed: MemberCount")
		} else if ring1.Weight != ringWeights[index1] || ring2.Weight != ringWeights[index2] {
			t.Fatal("SaveCoin failed: Weight")
		} else if ring1.Investor != ringInvestors[index1] || ring2.Investor != ringInvestors[index2] {
			t.Fatal("SaveCoin failed: Investor")
		} else if ring1.Rounds != pkg.RoundsCount || ring2.Rounds != pkg.RoundsCount {
			t.Fatal("SaveCoin failed: Rounds")
		} else if index1^index2 != 7 {
			t.Fatal("SaveCoin failed: Cooperation ring not match")
		}

		saveAndCheck(trader, coin7, false, false)
	}
}
