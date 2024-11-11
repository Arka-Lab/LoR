package internal

import (
	"errors"
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
	Locker   sync.Mutex
	Traders  map[string]*pkg.Trader
	Coins    map[string]pkg.CoinTable
	Fractals map[string]*pkg.FractalRing
}

func NewSystem() *System {
	return &System{
		Locker:   sync.Mutex{},
		Traders:  make(map[string]*pkg.Trader),
		Coins:    make(map[string]pkg.CoinTable),
		Fractals: make(map[string]*pkg.FractalRing),
	}
}

func (system *System) ProcessCoin(coin pkg.CoinTable) error {
	trader := system.Traders[coin.Owner]
	if trader == nil {
		return errors.New("trader not found")
	}

	system.Locker.Lock()
	system.Coins[coin.ID] = coin
	fractal, err := trader.SaveCoin(coin, true)
	if err != nil {
		system.Locker.Unlock()
		return err
	}

	for _, t := range system.Traders {
		if t.ID != trader.ID {
			if _, err := t.SaveCoin(coin, false); err != nil {
				system.Locker.Unlock()
				return err
			}
		}
	}

	if fractal == nil {
		fractal = system.checkForAnotherFractal(trader.ID)
		if fractal == nil {
			system.Locker.Unlock()
			return nil
		}
	}

	if err := system.processFractal(trader, fractal); err != nil {
		system.Locker.Unlock()
		return err
	}
	system.Locker.Unlock()

	return nil
}

func (system *System) checkForAnotherFractal(traderID string) *pkg.FractalRing {
	traderIDs := make([]string, 0)
	for id := range system.Traders {
		if id == traderID {
			traderIDs = append(traderIDs, id)
		}
	}
	rand.Shuffle(len(traderIDs), func(i, j int) {
		traderIDs[i], traderIDs[j] = traderIDs[j], traderIDs[i]
	})

	for _, id := range traderIDs {
		trader := system.Traders[id]
		if fractal := trader.CheckForRings(); fractal != nil {
			return fractal
		}
	}
	return nil
}

func (system *System) processFractal(trader *pkg.Trader, fractal *pkg.FractalRing) error {
	if err := system.verifyFractal(fractal); err != nil {
		trader.RemoveFractalRing(fractal.ID)
		return err
	}
	if err := system.blockCoins(fractal); err != nil {
		return err
	}
	if err := system.informOthers(fractal); err != nil {
		return err
	}
	system.Fractals[fractal.ID] = fractal
	return nil
}

func (system *System) verifyFractal(fractal *pkg.FractalRing) error {
	for _, traderID := range fractal.VerificationTeam {
		if trader, ok := system.Traders[traderID]; !ok {
			return errors.New("trader not found")
		} else if err := trader.SubmitRing(fractal); err != nil {
			return err
		}
	}
	return nil
}

func (system *System) blockCoins(fractal *pkg.FractalRing) error {
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

func (system *System) CreateRandomCoins(trader *pkg.Trader, done <-chan bool, errors chan<- error) {
	for {
		select {
		case <-done:
			return
		case <-trader.Data.Ticker.C:
			amount := rand.Float64() * 100
			coinType := rand.Intn(int(trader.Data.CoinTypeCount))
			if coin := trader.CreateCoin(amount, uint(coinType)); coin != nil {
				if err := system.ProcessCoin(*coin); err != nil {
					errors <- err
				}
			}
		}
	}
}

func (system *System) Init(numTraders int, coinTypeCount uint) error {
	ch := make(chan bool)
	for i := 0; i < numTraders; i++ {
		go func() {
			amount, wallet := rand.Float64()*100, uuid.New().String()
			trader := pkg.CreateTrader(amount, wallet, coinTypeCount)
			system.Traders[trader.ID] = trader
			ch <- true
		}()
	}
	for i := 0; i < numTraders; i++ {
		<-ch
	}

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
				syscall.Exit(0)
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
