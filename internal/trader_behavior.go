package internal

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"

	"github.com/Arka-Lab/LoR/pkg"
)

const (
	Debug = false
)

func RandomBehavior(numTraders, numTypes, numCoins int) error {
	if !Debug {
		log.SetOutput(io.Discard)
	}
	fmt.Println("Starting random behavior simulation...")

	traders, err := createTraders(numTraders, numTypes)
	if err != nil {
		return err
	}
	fmt.Println("Traders created successfully.")

	coins, err := createCoins(numCoins, numTypes, traders)
	if err != nil {
		return err
	}
	fmt.Println("Coins created successfully.")

	rings, fractals, err := processCoins(numCoins, numTraders, traders, coins)
	if err != nil {
		return err
	}
	fmt.Println("Coins processed successfully.")

	fmt.Println("Trader statistics:")
	for i := 0; i < numTraders; i++ {
		fmt.Printf("\tTrader %d has %d cooperation rings and %d fractal teams.\n", i+1, len(rings[i]), len(fractals[i]))
	}

	fmt.Println("Random behavior simulation completed.")
	return nil
}

func createTraders(numTraders, numTypes int) ([]*pkg.Trader, error) {
	ch1 := make(chan *pkg.Trader)
	for i := 0; i < numTraders; i++ {
		go func() {
			trader := pkg.CreateTrader(rand.Float64()*1000, "wallet-"+strconv.Itoa(i+1), numTypes)
			ch1 <- trader
		}()
	}

	traders := make([]*pkg.Trader, 0, numTraders)
	for i := 0; i < numTraders; i++ {
		trader := <-ch1
		if trader == nil {
			return nil, fmt.Errorf("failed to create trader")
		}
		traders = append(traders, trader)
	}

	ch2 := make(chan error)
	for _, t1 := range traders {
		go func(t1 *pkg.Trader) {
			for _, t2 := range traders {
				if err := t1.SaveTrader(*t2); err != nil {
					ch2 <- fmt.Errorf("failed to save trader: %v", err)
				}
			}
			ch2 <- nil
		}(t1)
	}

	for i := 0; i < numTraders; i++ {
		if err := <-ch2; err != nil {
			return nil, err
		}
	}
	return traders, nil
}

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
	for j := 0; j < numTraders; j++ {
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
		}(traders[j], coin, j)
	}

	for j := 0; j < numTraders; j++ {
		if err := <-ch; err != nil {
			return err
		}
	}
	return nil
}
