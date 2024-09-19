package internal

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/Arka-Lab/LoR/pkg"
	"github.com/Arka-Lab/LoR/tools"
)

func createCoins(numCoins, numTypes int, traders []*pkg.Trader) ([]*pkg.CoinTable, error) {
	ch := make(chan *pkg.CoinTable)
	index := make([]int, numCoins)
	for i := 0; i < numCoins; i++ {
		go func(i int) {
			randomIndex, coinType := rand.Intn(len(traders)), uint(rand.Intn(numTypes+1))
			trader := traders[randomIndex]

			coin := trader.CreateCoin(rand.Float64()*1000, coinType)
			index[i] = randomIndex
			ch <- coin
		}(i)
	}

	coins := make([]*pkg.CoinTable, 0, numCoins)
	for i := 0; i < numCoins; i++ {
		coin := <-ch
		if coin == nil {
			return nil, fmt.Errorf("failed to create coin %d for trader %d", i+1, index[i]+1)
		}
		coins = append(coins, coin)
	}
	return coins, nil
}

func processCoins(numCoins, numTraders int, traders []*pkg.Trader, coins []*pkg.CoinTable) ([][]*pkg.CooperationTable, [][][]string, error) {
	rings := make([][]*pkg.CooperationTable, numTraders)
	fractals := make([][][]string, numTraders)
	for i := 0; i < numCoins; i++ {
		log.Printf("Processing coin %d...", i+1)
		if err := processCoinForAllTraders(i, numTraders, traders, coins[i], rings, fractals); err != nil {
			return nil, nil, err
		}
		log.Printf("Coin %d processed successfully.", i+1)
	}
	return rings, fractals, nil
}

func processCoinForAllTraders(coinIndex, numTraders int, traders []*pkg.Trader, coin *pkg.CoinTable, rings [][]*pkg.CooperationTable, fractals [][][]string) error {
	ch := make(chan error)
	for i := 0; i < numTraders; i++ {
		go func(trader *pkg.Trader, coin *pkg.CoinTable, traderIndex int) {
			ring, team, err := trader.SaveCoin(*coin)
			if err != nil {
				ch <- fmt.Errorf("failed to save coin %d to trader %d: %v", coinIndex+1, traderIndex+1, err)
				return
			}

			if ring != nil {
				rings[traderIndex] = append(rings[traderIndex], ring)
				if team != nil {
					fractals[traderIndex] = append(fractals[traderIndex], team)
				}
			}
			ch <- nil
		}(traders[i], coin, i)
	}

	for i := 0; i < numTraders; i++ {
		if err := <-ch; err != nil {
			return err
		}
	}
	return nil
}

func analyzeCoins(numTypes, numCoins int, traders []*pkg.Trader, coins []*pkg.CoinTable, rings [][]*pkg.CooperationTable) {
	mp := make(map[string]int)
	for index, traderRings := range rings {
		trader := traders[index]
		for _, ring := range traderRings {
			current := trader.Data.Coins[ring.Investor]
			for i := 0; i <= numTypes; i++ {
				mp[current.ID]++
				current = trader.Data.Coins[ring.Next]
			}
		}
	}

	coinIDs := make([]string, 0, numCoins)
	for _, coin := range coins {
		coinIDs = append(coinIDs, coin.ID)
	}

	fmt.Println("Coin statistics:")
	for index, coinID := range tools.RandomSelect(coinIDs, 10) {
		fmt.Printf("\tCoin %d has %d copies.\n", index+1, mp[coinID])
	}
}
