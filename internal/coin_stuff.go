package internal

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/Arka-Lab/LoR/pkg"
)

func createCoins(numCoins, numTypes int, traders []*pkg.Trader) ([]*pkg.CoinTable, error) {
	ch := make(chan bool)
	coins := make([]*pkg.CoinTable, numCoins)
	for i := 0; i < numCoins; i++ {
		go func(i int) {
			trader := traders[rand.Intn(len(traders))]
			coinType := uint(rand.Intn(numTypes + 1))

			coins[i] = trader.CreateCoin(rand.Float64()*1000, coinType)
			ch <- coins[i] != nil
		}(i)
	}

	for i := 0; i < numCoins; i++ {
		if ok := <-ch; !ok {
			return nil, fmt.Errorf("failed to create coin")
		}
	}
	return coins, nil
}

func processCoins(numCoins int, traders []*pkg.Trader, coins []*pkg.CoinTable) ([][]*pkg.CooperationTable, [][][]string, error) {
	rings := make([][]*pkg.CooperationTable, len(traders))
	fractals := make([][][]string, len(traders))
	for i := 0; i < numCoins; i++ {
		log.Printf("Processing coin %d...", i+1)
		if ring, fractal, err := processCoinForTraders(traders, coins[i]); err != nil {
			return nil, nil, err
		} else {
			for j := 0; j < len(traders); j++ {
				if ring[j] != nil {
					rings[j] = append(rings[j], ring[j])
				}
				if fractal[j] != nil {
					fractals[j] = append(fractals[j], fractal[j])
				}
			}
		}
		log.Printf("Coin %d processed successfully.", i+1)
	}
	return rings, fractals, nil
}

func processCoinForTraders(traders []*pkg.Trader, coin *pkg.CoinTable) ([]*pkg.CooperationTable, [][]string, error) {
	ch := make(chan error)
	fractals := make([][]string, len(traders))
	rings := make([]*pkg.CooperationTable, len(traders))
	for i := 0; i < len(traders); i++ {
		go func(trader *pkg.Trader, coin *pkg.CoinTable, traderIndex int) {
			ring, team, err := trader.SaveCoin(*coin)
			if err != nil {
				ch <- fmt.Errorf("failed to save coin to trader %d: %v", traderIndex+1, err)
				return
			}

			rings[traderIndex], fractals[traderIndex] = ring, team
			ch <- nil
		}(traders[i], coin, i)
	}

	for i := 0; i < len(traders); i++ {
		if err := <-ch; err != nil {
			return nil, nil, err
		}
	}
	return rings, fractals, nil
}

func analyzeCoins(numTypes int, traders []*pkg.Trader, rings [][]*pkg.CooperationTable) {
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
	for i := len(traders) - 1; i >= 0; i-- {
		counter[i] += counter[i+1]
	}

	fmt.Println("Coin statistics:")
	for i := 0; i <= len(traders); i++ {
		fmt.Printf("\tCoins with at least %d traders: %d\n", i, counter[i])
	}
}
