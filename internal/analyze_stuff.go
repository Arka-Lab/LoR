package internal

import (
	"fmt"
	"math/rand"
	"slices"

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

func analyzeFractals(rings [][]*pkg.CooperationTable, fractals [][]*pkg.FractalRing) {
	coins := make(map[string][]string)
	for i, traderFractals := range fractals {
		cooperations := make(map[string]*pkg.CooperationTable)
		for _, ring := range rings[i] {
			cooperations[ring.ID] = ring
		}

		for _, fractal := range traderFractals {
			for _, ringID := range fractal.CooperationRings {
				ring := cooperations[ringID]
				coins[fractal.ID] = append(coins[fractal.ID], ring.Members...)
			}
			slices.Sort(coins[fractal.ID])
		}
	}

	hasFractal := make([]int, 0)
	for i, traderFractals := range fractals {
		if len(traderFractals) > 0 {
			hasFractal = append(hasFractal, i)
		}
	}

	T := 1000000
	fmt.Println("Fractal statistics:")
	if len(hasFractal) == 0 {
		fmt.Println("\tNo fractal rings found.")
		return
	}
	for i, prev, total := 1, 1, 0; i <= T; i++ {
		index1 := hasFractal[rand.Intn(len(hasFractal))]
		idx1 := rand.Intn(len(fractals[index1]))
		fractal1 := fractals[index1][idx1]

		index2 := hasFractal[rand.Intn(len(hasFractal))]
		idx2 := rand.Intn(len(fractals[index2]))
		fractal2 := fractals[index2][idx2]

		commons := 0
		for _, coinID := range coins[fractal1.ID] {
			if _, ok := slices.BinarySearch(coins[fractal2.ID], coinID); ok {
				commons++
			}
		}
		total += commons

		if i == prev {
			fmt.Printf("\tAverage number of common coins between two fractal rings after %d iterations: %.2f\n", i, float64(total)/float64(i))
			prev *= 10
		}
	}
}
