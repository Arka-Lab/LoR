package pkg_test

import (
	"fmt"
	"math/rand"
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

	if coin := trader.CreateCoin(100, 1); tools.VerifyWithPublicKeyStr(trader.ID+"-1", coin.ID, trader.PublicKey) != nil {
		t.Error("CreateCoin failed: first ID")
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

	if coin := trader.CreateCoin(96.5, 1); tools.VerifyWithPublicKeyStr(trader.ID+"-1", coin.ID, trader.PublicKey) != nil {
		t.Error("CreateCoin failed: second ID")
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
	trader := pkg.CreateTrader(10.5, "test_wallet", 2)
	coin1 := trader.CreateCoin(10.5, 3)
	coin2 := trader.CreateCoin(13.8, 2)
	trader.SaveTrader(*trader)

	saveCoinAndCheck(t, trader, coin1, true, false, false)
	saveCoinAndCheck(t, trader, coin2, false, false, false)
	saveCoinAndCheck(t, trader, coin2, true, false, false)
}

func TestCooperationRing(t *testing.T) {
	traders := createTrader(2, 2)
	coins := []*pkg.CoinTable{
		traders[0].CreateCoin(10.5, 0),
		traders[1].CreateCoin(9.8, 1),
		traders[1].CreateCoin(7.3, 0),
		traders[1].CreateCoin(2.5, 1),
		traders[1].CreateCoin(3.4, 2),
		traders[0].CreateCoin(7.3, 2),
		traders[0].CreateCoin(11.7, 2),
	}

	ringIDs, ringWeights, ringInvestors := saveCoins([]*pkg.CoinTable{coins[0], coins[2]}, []*pkg.CoinTable{coins[1], coins[3]}, []*pkg.CoinTable{coins[4], coins[5]})
	for _, trader := range traders {
		saveBatch(t, trader, coins[:4], false, false, false)
		rings, _ := saveBatch(t, trader, coins[4:6], false, false, true)
		ring1, ring2 := rings[0], rings[1]
		if ring1.ID == ring2.ID {
			t.Fatal("SaveCoin failed: Same cooperation ring ID")
		}
		testRing(t, trader, ring1)
		testRing(t, trader, ring2)

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

		saveCoinAndCheck(t, trader, coins[6], false, false, false)
	}
}

func TestFractalRing(t *testing.T) {
	testcase := []struct {
		traderCount int
		hasTeam     bool
	}{{pkg.VerificationMax << 1, true}, {pkg.VerificationMin - 1, false}}

	for _, tc := range testcase {
		var team []string
		traders := createTrader(tc.traderCount, 2)
		for i := 0; team == nil && i < pkg.FractalMax; i++ {
			coin1 := traders[rand.Intn(len(traders))].CreateCoin(100, 0)
			coin2 := traders[rand.Intn(len(traders))].CreateCoin(100, 1)
			coin3 := traders[rand.Intn(len(traders))].CreateCoin(100, 2)
			for _, trader := range traders {
				saveBatch(t, trader, []*pkg.CoinTable{coin1, coin2}, false, false, false)
				_, team = saveCoinAndCheck(t, trader, coin3, false, tc.hasTeam, true)
			}
		}

		if team == nil && tc.hasTeam {
			t.Fatal("SaveCoin failed: Fractal ring not found")
		} else if team != nil && !tc.hasTeam {
			t.Fatal("SaveCoin failed: Fail to detect fractal ring")
		}
	}
}

func createTrader(count, types int) []*pkg.Trader {
	traders := make([]*pkg.Trader, count)
	for i := 0; i < count; i++ {
		traders[i] = pkg.CreateTrader(float64(i+1)*10, fmt.Sprintf("test_wallet%d", i), types)
	}

	for _, trader1 := range traders {
		for _, trader2 := range traders {
			trader1.SaveTrader(*trader2)
		}
	}
	return traders
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
				ringIDs = append(ringIDs, tools.SHA256Str([]string{c1.ID, c2.ID, c3.ID}))
				ringWeights = append(ringWeights, c1.Amount+c2.Amount+c3.Amount)
				ringInvestors = append(ringInvestors, c1.ID)
			}
		}
	}
	return
}

func testRing(t *testing.T, trader *pkg.Trader, ring *pkg.CooperationTable) {
	coins := trader.Data.Coins
	for i, current := 0, ring.Members[0]; i < len(ring.Members); i++ {
		coin := coins[current]
		if coin.Status != pkg.Blocked {
			t.Fatal("SaveCoin failed: Coin status")
		} else if coin.Prev != ring.Members[(i-1+len(ring.Members))%len(ring.Members)] {
			t.Fatal("SaveCoin failed: Prev")
		} else if coin.Next != ring.Members[(i+1)%len(ring.Members)] {
			t.Fatal("SaveCoin failed: Next")
		}
		current = coin.Next
	}

	if ring.ID != tools.SHA256Str(ring.Members) {
		t.Fatal("SaveCoin failed: Cooperation ring ID")
	} else if ring.Investor != ring.Members[0] {
		t.Fatal("SaveCoin failed: Investor")
	} else if ring.Weight != coins[ring.Members[0]].Amount+coins[ring.Members[1]].Amount+coins[ring.Members[2]].Amount {
		t.Fatal("SaveCoin failed: Weight")
	} else if ring.Rounds != pkg.RoundsCount {
		t.Fatal("SaveCoin failed: Rounds")
	}
}
