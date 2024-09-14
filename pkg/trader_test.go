package pkg_test

import (
	"testing"

	"github.com/Arka-Lab/LoR/pkg"
	"github.com/Arka-Lab/LoR/tools"
)

func TestCreateTrader(t *testing.T) {
	trader := pkg.CreateTrader(100, "test_wallet", 10)

	if trader.ID != tools.SHA256Str("test_wallet-10") {
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

	if coin := trader.CreateCoin(100, 1); coin.ID != tools.SHA256Str(trader.ID+"-1-1") {
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

	if coin := trader.CreateCoin(96.5, 1); coin.ID != tools.SHA256Str(trader.ID+"-1-2") {
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

func saveCoinAndCheck(t *testing.T, trader *pkg.Trader, coin *pkg.CoinTable, hasError, hasFractal, hasRing bool) (*pkg.CooperationTable, []string) {
	ring, fractal, err := trader.SaveCoin(*coin)
	if err != nil && !hasError {
		t.Fatal("SaveCoin failed:", err)
	} else if err == nil && hasError {
		t.Fatal("SaveCoin failed: Error not detected")
	} else if ring != nil && !hasRing {
		t.Fatal("SaveCoin failed: Cooperation ring detected")
	} else if ring == nil && hasRing {
		t.Fatal("SaveCoin failed: Fail to detect cooperation ring")
	} else if fractal != nil && !hasFractal {
		t.Fatal("SaveCoin failed: Fractal ring detected")
	} else if fractal == nil && hasFractal {
		t.Fatal("SaveCoin failed: Fail to detect fractal ring")
	}
	return ring, fractal
}

func saveBatch(t *testing.T, trader *pkg.Trader, coins []*pkg.CoinTable, hasError, hasFractal, hasRing bool) (rings []*pkg.CooperationTable, fractals [][]string) {
	for _, coin := range coins {
		ring, fractal := saveCoinAndCheck(t, trader, coin, hasError, hasFractal, hasRing)
		rings, fractals = append(rings, ring), append(fractals, fractal)
	}
	return
}

func TestSaveCoin(t *testing.T) {
	trader := pkg.CreateTrader(10.5, "test_wallet", 2)
	coin1 := trader.CreateCoin(10.5, 3)
	coin2 := trader.CreateCoin(13.8, 2)

	saveCoinAndCheck(t, trader, coin1, true, false, false)
	saveCoinAndCheck(t, trader, coin2, false, false, false)
	saveCoinAndCheck(t, trader, coin2, true, false, false)
}

func findIndex(arr []string, val string) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
}

func saveCoins(c1s, c2s, c3s []*pkg.CoinTable) (ringIDs []string, ringWeights []float64, ringInvestors []string) {
	for _, c1 := range c1s {
		for _, c2 := range c2s {
			for _, c3 := range c3s {
				ringIDs = append(ringIDs, tools.SHA256Str(c1.ID+"-"+c2.ID+"-"+c3.ID))
				ringWeights = append(ringWeights, c1.Amount+c2.Amount+c3.Amount)
				ringInvestors = append(ringInvestors, c1.ID)
			}
		}
	}
	return
}

func TestCooperationRing(t *testing.T) {
	trader1 := pkg.CreateTrader(100, "test_wallet1", 2)
	trader2 := pkg.CreateTrader(100, "test_wallet2", 2)

	coin1 := trader1.CreateCoin(10.5, 0)
	coin2 := trader2.CreateCoin(9.8, 1)
	coin3 := trader2.CreateCoin(7.3, 0)
	coin4 := trader2.CreateCoin(2.5, 1)
	coin5 := trader2.CreateCoin(3.4, 2)
	coin6 := trader1.CreateCoin(7.3, 2)
	coin7 := trader1.CreateCoin(11.7, 2)

	ringIDs, ringWeights, ringInvestors := saveCoins([]*pkg.CoinTable{coin1, coin3}, []*pkg.CoinTable{coin2, coin4}, []*pkg.CoinTable{coin5, coin6})
	for _, trader := range []*pkg.Trader{trader1, trader2} {
		saveBatch(t, trader, []*pkg.CoinTable{coin1, coin2, coin3, coin4}, false, false, false)
		rings, _ := saveBatch(t, trader, []*pkg.CoinTable{coin5, coin6}, false, false, true)
		ring1, ring2 := rings[0], rings[1]
		if ring1.ID == ring2.ID {
			t.Fatal("SaveCoin failed: Same cooperation ring ID")
		}

		index1, index2 := findIndex(ringIDs, ring1.ID), findIndex(ringIDs, ring2.ID)
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

		saveCoinAndCheck(t, trader, coin7, false, false, false)
	}
}
