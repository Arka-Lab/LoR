package internal

import (
	"fmt"

	"github.com/Arka-Lab/LoR/pkg"
)

func analyzeTraders(numTraders int, rings [][]*pkg.CooperationTable, fractals [][]*pkg.FractalRing) {
	fmt.Println("Trader statistics:")
	for i := 0; i < numTraders; i++ {
		fmt.Printf("\tTrader %d has %d cooperation rings and %d fractal rings.\n", i+1, len(rings[i]), len(fractals[i]))
	}
}

func analyzeRings(numTypes int, traders []*pkg.Trader, rings [][]*pkg.CooperationTable) {
	mp := make(map[string]int)
	for index, traderRings := range rings {
		coins := traders[index].Data.Coins
		for _, ring := range traderRings {
			for i, current := 0, ring.Investor; i <= numTypes; i++ {
				mp[current]++
				current = coins[current].Next
			}
		}
	}

	counter := make([]int, len(traders)+1)
	for _, count := range mp {
		counter[count]++
	}
	for i := len(traders) - 1; i > 0; i-- {
		counter[i] += counter[i+1]
	}

	fmt.Println("Coin statistics:")
	for i := 1; i <= len(traders); i++ {
		fmt.Printf("\tCoins in at least %d cooperation rings: %d\n", i, counter[i])
	}
}
