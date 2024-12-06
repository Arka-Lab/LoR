package internal

import (
	"fmt"

	"github.com/Arka-Lab/LoR/pkg"
)

func AnalyzeSystem(system *System) {
	fmt.Println("Number of coins:", len(system.Coins))
	fmt.Println("Number of fractal rings:", len(system.Fractals))

	runCoins := 0
	for _, coin := range system.Coins {
		if coin.Status == pkg.Run {
			runCoins++
		}
	}
	fmt.Printf("Percentage of run coins: %.2f%%\n", float64(runCoins)/float64(len(system.Coins))*100)

	numSubmitted, totalSubmitted, acceptRate := 0, 0, 0.0
	for traderID := range system.Traders {
		if system.SubmitCount[traderID] > 0 {
			numSubmitted++
			totalSubmitted += system.SubmitCount[traderID]
			acceptRate += float64(system.AcceptedCount[traderID]) / float64(system.SubmitCount[traderID])
		}
	}
	fmt.Printf("Average number of submitted fractal rings per trader: %.2f\n", float64(totalSubmitted)/float64(numSubmitted))
	fmt.Printf("Average fractal ring acceptance rate per trader: %.2f%%\n", acceptRate/float64(numSubmitted)*100)

	fmt.Println("Number of invalid accepted fractal rings:", system.BadAcceptCount)
	fmt.Println("Number of valid rejected fractal rings:", system.BadRejectCount)

	if RunFractals {
		coinsCount, coinsTotal := 0, 0.
		coinsSatisfaction := make(map[string]float64)
		for _, fractal := range system.Fractals {
			for _, ring := range fractal.CooperationRings {
				if ring.Rounds != -1 {
					satisfaction := float64(ring.Rounds) / float64(pkg.RoundsCount)
					if !ring.IsValid {
						satisfaction *= -1
					}

					coinsCount += len(ring.CoinIDs)
					coinsTotal += satisfaction * float64(len(ring.CoinIDs))
					for _, coinID := range ring.CoinIDs {
						coinsSatisfaction[coinID] = satisfaction
					}
				}
			}
		}
		fmt.Printf("Average satisfaction per coin: %.2f%%\n", float64(coinsTotal)/float64(coinsCount)*100)

		traderSatisfaction := make(map[string][]float64)
		for coinID, satisfaction := range coinsSatisfaction {
			owner := system.Coins[coinID].Owner
			traderSatisfaction[owner] = append(traderSatisfaction[owner], satisfaction)
		}

		tradersTotal := 0.
		for _, satisfactions := range traderSatisfaction {
			total := 0.
			for _, satisfaction := range satisfactions {
				total += satisfaction
			}
			tradersTotal += total / float64(len(satisfactions))
		}
		fmt.Printf("Average satisfaction per trader: %.2f%%\n", float64(tradersTotal)/float64(len(traderSatisfaction))*100)

		adjacentFractals := make(map[string]map[string]bool)
		tradersAdjacency := make(map[string]map[string]bool)
		verificationAdjacency := make(map[string]map[string]bool)
		for _, trader := range system.Traders {
			adjacentFractals[trader.ID] = make(map[string]bool)
			tradersAdjacency[trader.ID] = make(map[string]bool)
			verificationAdjacency[trader.ID] = make(map[string]bool)
		}
		for _, fractal := range system.Fractals {
			for _, ring := range fractal.CooperationRings {
				for _, coinID := range ring.CoinIDs {
					owner := system.Coins[coinID].Owner
					adjacentFractals[owner][fractal.ID] = true

					for _, traderID := range fractal.VerificationTeam {
						tradersAdjacency[owner][traderID] = true
						verificationAdjacency[traderID][owner] = true
					}
					for _, otherCoinID := range ring.CoinIDs {
						tradersAdjacency[owner][system.Coins[otherCoinID].Owner] = true
					}
				}
			}
		}

		tradersCount := 0
		totalAdjacency, maximumAdjacency := 0, 0
		for traderID := range system.Traders {
			adjacency := len(verificationAdjacency[traderID]) / len(system.Fractals)
			if len(adjacentFractals[traderID]) > 0 {
				adjacency += len(tradersAdjacency[traderID]) / len(adjacentFractals[traderID])
			}

			if adjacency > 0 {
				tradersCount++
				totalAdjacency += adjacency
				if adjacency > maximumAdjacency {
					maximumAdjacency = adjacency
				}
			}
		}
		fmt.Printf("Average adjacency per trader: %.2f\n", float64(totalAdjacency)/float64(tradersCount))
		fmt.Printf("Maximum adjacency per trader: %d\n", maximumAdjacency)
	}
}
