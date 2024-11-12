package internal

import (
	"fmt"

	"github.com/Arka-Lab/LoR/pkg"
)

func AnalyzeSystem(system *System) {
	fmt.Println("Number of coins:", len(system.Coins))
	fmt.Println("Number of fractal rings:", len(system.Fractals))

	numCoins := 0
	for _, coin := range system.Coins {
		if coin.Status == pkg.Blocked {
			numCoins++
		}
	}
	fmt.Printf("Percentage of blocked coins: %.2f%%\n", float64(numCoins)/float64(len(system.Coins))*100)

	numSubmitted, totalSubmitted, acceptRate := 0, 0, 0.0
	for traderID := range system.Traders {
		if system.SubmitCount[traderID] > 0 {
			numSubmitted++
			totalSubmitted += system.SubmitCount[traderID]
			acceptRate += float64(system.AcceptedCount[traderID]) / float64(system.SubmitCount[traderID])
		}
	}
	fmt.Printf("Average number of submitted rings per trader: %.2f\n", float64(totalSubmitted)/float64(numSubmitted))
	fmt.Printf("Average acceptance rate per trader: %.2f%%\n", acceptRate/float64(numSubmitted)*100)
}
