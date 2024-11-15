package internal

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"syscall"

	"github.com/Arka-Lab/LoR/pkg"
	"github.com/google/uuid"
)

const (
	Debug = true
)

type System struct {
	Locker        sync.Mutex
	SubmitCount   map[string]int
	AcceptedCount map[string]int
	Traders       map[string]*pkg.Trader
	Coins         map[string]pkg.CoinTable
	Fractals      map[string]*pkg.FractalRing
}

func NewSystem() *System {
	return &System{
		Locker:        sync.Mutex{},
		SubmitCount:   make(map[string]int),
		AcceptedCount: make(map[string]int),
		Traders:       make(map[string]*pkg.Trader),
		Coins:         make(map[string]pkg.CoinTable),
		Fractals:      make(map[string]*pkg.FractalRing),
	}
}

func (system *System) ProcessCoin(coin pkg.CoinTable) error {
	system.Locker.Lock()
	defer system.Locker.Unlock()

	system.Coins[coin.ID] = coin
	for _, t := range system.Traders {
		if err := t.SaveCoin(coin); err != nil {
			return err
		}
	}

	for index, traderID := range system.getShuffledTraderIDs(coin.Owner) {
		trader := system.Traders[traderID]
		if fractal := trader.CheckForRings(); fractal != nil {
			system.SubmitCount[traderID]++
			if err := system.processFractal(trader, fractal); err != nil {
				return err
			}

			if Debug {
				log.Printf("Fractal ring created by trader %d with %d cooperation rings and %d verification team members\n", index+1, len(fractal.CooperationRings), len(fractal.VerificationTeam))
			}
			if err := system.runFractal(fractal); err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

func (system *System) getShuffledTraderIDs(firstID string) (result []string) {
	for traderID := range system.Traders {
		if traderID != firstID {
			result = append(result, traderID)
		}
	}
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	if firstID == "" {
		result = append([]string{firstID}, result...)
	}
	return
}

func (system *System) processFractal(trader *pkg.Trader, fractal *pkg.FractalRing) error {
	if err := system.verifyFractal(fractal); err != nil {
		trader.RemoveFractalRing(fractal.ID)
		return err
	} else if err := system.checkCoins(fractal); err != nil {
		return err
	} else if err := system.informOthers(fractal); err != nil {
		return err
	}
	system.Fractals[fractal.ID] = fractal
	system.AcceptedCount[trader.ID]++
	return nil
}

func (system *System) verifyFractal(fractal *pkg.FractalRing) error {
	totalErrors := []error{}
	for _, traderID := range fractal.VerificationTeam {
		if trader, ok := system.Traders[traderID]; !ok {
			return errors.New("trader not found")
		} else if err := trader.SubmitRing(fractal); err != nil {
			totalErrors = append(totalErrors, err)
		}
	}

	if 2*len(totalErrors) >= len(fractal.VerificationTeam) {
		fmt.Println(len(totalErrors), len(fractal.VerificationTeam))
		fmt.Println(totalErrors[0])
		return errors.New("fractal ring verification failed")
	}
	return nil
}

func (system *System) checkCoins(fractal *pkg.FractalRing) error {
	for _, ring := range fractal.CooperationRings {
		for _, coinID := range ring.CoinIDs {
			if coin, ok := system.Coins[coinID]; !ok {
				return errors.New("coin not found")
			} else if coin.Status != pkg.Run {
				return errors.New("coin is not running")
			}
		}
	}
	return nil
}

func (system *System) informOthers(fractal *pkg.FractalRing) error {
	for _, ring := range fractal.CooperationRings {
		for _, coinID := range ring.CoinIDs {
			coin := system.Coins[coinID]
			coin.Status = pkg.Blocked
			system.Coins[coinID] = coin
		}
	}
	for _, trader := range system.Traders {
		if err := trader.InformFractalRing(*fractal); err != nil {
			return err
		}
	}
	return nil
}

func (system *System) runFractal(fractal *pkg.FractalRing) error {
	for round := 0; round < pkg.RoundsCount; round++ {
		for _, traderID := range fractal.VerificationTeam {
			if trader, ok := system.Traders[traderID]; !ok {
				return errors.New("trader not found")
			} else if err := trader.Vote(); err != nil {
				if err2 := system.applyRounds(fractal, round); err2 != nil {
					return err2
				}
				return err
			}
		}
	}

	return system.applyRounds(fractal, pkg.RoundsCount)
}

func (system *System) applyRounds(fractal *pkg.FractalRing, round int) error {
	fractal.SuccessfulRounds = round
	for _, ring := range fractal.CooperationRings {
		money := system.Coins[ring.CoinIDs[0]].Amount * float64(round) / pkg.RoundsCount
		if err := system.applyRing(ring, money, round); err != nil {
			return err
		}
	}
	return nil
}

func (system *System) applyRing(ring pkg.CooperationTable, money float64, round int) error {
	for _, coinID := range ring.CoinIDs {
		coin := system.Coins[coinID]
		amount := money * coin.Amount / ring.Weight
		if round < pkg.RoundsCount {
			coin.Status = pkg.Expired
		} else {
			coin.Status = pkg.Paid
			amount += pkg.FractalPrize
		}
		system.Coins[coinID] = coin

		for _, trader := range system.Traders {
			if err := trader.UpdateBalance(coin.Owner, amount); err != nil {
				return err
			}
		}
	}

	for _, trader := range system.Traders {
		if round < pkg.RoundsCount {
			trader.ExpireRing(ring)
		} else {
			trader.PayRing(ring)
		}
	}
	return nil
}

func (system *System) CreateRandomCoins(trader *pkg.Trader, done <-chan bool, errors chan<- error) {
	for {
		select {
		case <-done:
			return
		case <-trader.Data.Ticker.C:
			amount := rand.Float64() * 10
			coinType := rand.Intn(int(trader.Data.CoinTypeCount))
			if coin := trader.CreateCoin(amount, uint(coinType)); coin != nil {
				if err := system.ProcessCoin(*coin); err != nil {
					errors <- err
				}
			}
		}
	}
}

func (system *System) Init(numTraders, numRandomVoters, numBadVoters int, coinTypeCount uint) error {
	ch := make(chan bool)
	for i := 0; i < numTraders; i++ {
		go func() {
			var trader *pkg.Trader
			amount, wallet := rand.Float64()*1000, uuid.New().String()
			if i < numRandomVoters {
				trader = pkg.CreateTrader(pkg.RandomVote, amount, wallet, coinTypeCount)
			} else if i < numRandomVoters+numBadVoters {
				trader = pkg.CreateTrader(pkg.BadVote, amount, wallet, coinTypeCount)
			} else {
				trader = pkg.CreateTrader(pkg.Normal, amount, wallet, coinTypeCount)
			}

			system.Locker.Lock()
			system.Traders[trader.ID] = trader
			system.Locker.Unlock()
			ch <- true
		}()
	}
	log.Printf("%d traders created: %d random voters, %d bad voters\n", numTraders, numRandomVoters, numBadVoters)
	for i := 0; i < numTraders; i++ {
		<-ch
	}
	return system.saveTraders()
}

func (system *System) saveTraders() error {
	for _, trader1 := range system.Traders {
		for _, trader2 := range system.Traders {
			if err := trader1.SaveTrader(*trader2); err != nil {
				return err
			}
		}
	}
	return nil
}

func (system *System) Start(finish <-chan bool) {
	errors := make(chan error)
	dones := make([]chan bool, len(system.Traders))

	i := 0
	for _, trader := range system.Traders {
		dones[i] = make(chan bool, 1)
		go system.CreateRandomCoins(trader, dones[i], errors)
		i++
	}

	finished := false
	for !finished {
		select {
		case err := <-errors:
			if Debug {
				log.Println("Error:", err)
				if err.Error() != "bad behavior" {
					syscall.Exit(1)
				}
			}
		case <-finish:
			i := 0
			for _, trader := range system.Traders {
				dones[i] <- true
				trader.Data.Ticker.Stop()
				i++
			}
			finished = true
		}
	}
}
