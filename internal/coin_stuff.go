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

func processCoins(numCoins int, traders []*pkg.Trader, coins []*pkg.CoinTable) ([][]*pkg.CooperationTable, [][]*pkg.FractalRing, error) {
	rings := make([][]*pkg.CooperationTable, len(traders))
	fractals := make([][]*pkg.FractalRing, len(traders))
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

func processCoinForTraders(traders []*pkg.Trader, coin *pkg.CoinTable) ([]*pkg.CooperationTable, []*pkg.FractalRing, error) {
	ch := make(chan error)
	fractals := make([]*pkg.FractalRing, len(traders))
	rings := make([]*pkg.CooperationTable, len(traders))
	for i := 0; i < len(traders); i++ {
		go func(trader *pkg.Trader, coin *pkg.CoinTable, traderIndex int) {
			ring, fractal, err := trader.SaveCoin(*coin)
			if err != nil {
				ch <- fmt.Errorf("failed to save coin to trader %d: %v", traderIndex+1, err)
				return
			}

			rings[traderIndex], fractals[traderIndex] = ring, fractal
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
