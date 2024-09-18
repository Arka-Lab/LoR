package internal

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/Arka-Lab/LoR/pkg"
)

const (
	Debug = false
)

func RandomBehavior(numTraders, numTypes, numCoins int) error {
	if !Debug {
		log.SetOutput(io.Discard)
	}
	log.Print("Starting random behavior simulation...")

	traders, err := createTraders(numTraders, numTypes)
	if err != nil {
		return err
	}
	log.Print("Traders created successfully.")

	coins, err := createCoins(numCoins, numTypes, traders)
	if err != nil {
		return err
	}
	log.Print("Coins created successfully.")

	rings, fractals, err := processCoins(numCoins, numTraders, traders, coins)
	if err != nil {
		return err
	}

	for i := 0; i < numTraders; i++ {
		fmt.Printf("Trader %d has %d cooperation rings and %d fractal teams.\n", i+1, len(rings[i]), len(fractals[i]))
	}

	log.Print("Random behavior simulation completed.")
	return nil
}

func createTraders(numTraders, numTypes int) ([]*pkg.Trader, error) {
	traders := make([]*pkg.Trader, 0, numTraders)
	for i := 0; i < numTraders; i++ {
		trader := pkg.CreateTrader(rand.Float64()*1000, "wallet-"+strconv.Itoa(i+1), numTypes)
		if trader == nil {
			return nil, fmt.Errorf("failed to create trader %d", i+1)
		}
		traders = append(traders, trader)
	}

	for index1, trader1 := range traders {
		for index2, trader2 := range traders {
			if err := trader1.SaveTrader(*trader2); err != nil {
				return nil, fmt.Errorf("failed to save trader %d to trader %d: %v", index1, index2, err)
			}
		}
	}
	return traders, nil
}

func createCoins(numCoins, numTypes int, traders []*pkg.Trader) ([]*pkg.CoinTable, error) {
	coins := make([]*pkg.CoinTable, 0, numCoins)
	for i := 0; i < numCoins; i++ {
		randomIndex, coinType := rand.Intn(len(traders)), uint(rand.Intn(numTypes+1))

		trader := traders[randomIndex]
		coin := trader.CreateCoin(rand.Float64()*1000, coinType)
		if coin == nil {
			return nil, fmt.Errorf("failed to create coin %d for trader %d", i+1, randomIndex+1)
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
		for j := 0; j < numTraders; j++ {
			ring, team, err := traders[j].SaveCoin(*coins[i])
			if err != nil {
				return nil, nil, fmt.Errorf("failed to save coin %d to trader %d: %v", i+1, j+1, err)
			}

			if ring != nil {
				rings[j] = append(rings[j], ring)
				log.Printf("cooperation ring detected by trader %d: %s", j+1, arrToStr(findIndexCoins(coins, ring.Members)))

				if team != nil {
					fractals[j] = append(fractals[j], team)
					log.Printf("fractal team detected by trader %d: %s", j+1, arrToStr(findIndexTraders(traders, team)))
				}
			}
		}
		log.Printf("Coin %d processed successfully.", i+1)
	}
	return rings, fractals, nil
}

func findIndexTraders(traders []*pkg.Trader, IDs []string) (result []int) {
	indexes := make(map[string]int)
	for i, trader := range traders {
		indexes[trader.ID] = i
	}

	for _, ID := range IDs {
		if index, ok := indexes[ID]; ok {
			result = append(result, index)
		} else {
			log.Fatal("Trader ID not found")
		}
	}
	return
}

func findIndexCoins(coins []*pkg.CoinTable, IDs []string) (result []int) {
	indexes := make(map[string]int)
	for i, coin := range coins {
		indexes[coin.ID] = i
	}

	for _, ID := range IDs {
		if index, ok := indexes[ID]; ok {
			result = append(result, index)
		} else {
			log.Fatal("Coin ID not found")
		}
	}
	return
}

func arrToStr(arr []int) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ", "), "[]")
}
