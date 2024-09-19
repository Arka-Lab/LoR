package internal

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/Arka-Lab/LoR/pkg"
)

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

func analyzeTraders(numTraders int, rings [][]*pkg.CooperationTable, fractals [][][]string) {
	fmt.Println("Trader statistics:")
	for i := 0; i < numTraders; i++ {
		fmt.Printf("\tTrader %d has %d cooperation rings and %d fractal rings.\n", i+1, len(rings[i]), len(fractals[i]))
	}
}
