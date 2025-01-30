package internal

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/Arka-Lab/LoR/internal/models"
	"github.com/google/uuid"
)

const (
	Debug       = false
	RunFractals = true
)

type System struct {
	BadAcceptCount int
	BadRejectCount int
	Locker         sync.Mutex
	SubmitCount    map[string]int
	AcceptedCount  map[string]int
	Traders        map[string]*models.Trader
	Coins          map[string]models.CoinTable
	Fractals       map[string]*models.FractalRing
}

func NewSystem() *System {
	return &System{
		BadAcceptCount: 0,
		BadRejectCount: 0,
		Locker:         sync.Mutex{},
		SubmitCount:    make(map[string]int),
		AcceptedCount:  make(map[string]int),
		Traders:        make(map[string]*models.Trader),
		Coins:          make(map[string]models.CoinTable),
		Fractals:       make(map[string]*models.FractalRing),
	}
}

func (system *System) ProcessCoin(coin models.CoinTable) error {
	system.Locker.Lock()
	defer system.Locker.Unlock()

	system.Coins[coin.ID] = coin
	if err := system.saveCoinToTraders(coin); err != nil {
		return err
	}

	return system.processTradersForCoin(coin)
}

func (system *System) saveCoinToTraders(coin models.CoinTable) error {
	for _, t := range system.Traders {
		if err := t.SaveCoin(coin); err != nil {
			return err
		}
	}
	return nil
}

func (system *System) processTradersForCoin(coin models.CoinTable) error {
	for index, traderID := range system.getShuffledTraderIDs(coin.Owner) {
		trader := system.Traders[traderID]
		if fractal := trader.CheckForRings(); fractal != nil {
			system.SubmitCount[traderID]++
			if err := system.handleFractal(trader, fractal, index); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func (system *System) handleFractal(trader *models.Trader, fractal *models.FractalRing, index int) error {
	if err := system.processFractal(trader, fractal); err != nil {
		if fractal.IsValid {
			system.BadRejectCount++
		}
		return err
	}

	if !fractal.IsValid {
		system.BadAcceptCount++
	}
	if Debug {
		log.Printf("Fractal ring created by trader %d with %d cooperation rings and %d verification team members\n", index+1, len(fractal.CooperationRings), len(fractal.VerificationTeam))
	}
	if RunFractals {
		if err := system.runFractal(fractal); err != nil {
			return err
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

func (system *System) processFractal(trader *models.Trader, fractal *models.FractalRing) error {
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

func (system *System) verifyFractal(fractal *models.FractalRing) error {
	totalErrors := []error{}
	for _, traderID := range fractal.VerificationTeam {
		if err := system.Traders[traderID].SubmitRing(fractal); err != nil {
			totalErrors = append(totalErrors, err)
		}
	}

	if 2*len(totalErrors) >= len(fractal.VerificationTeam) {
		return errors.New("fractal ring verification failed")
	}
	return nil
}

func (system *System) checkCoins(fractal *models.FractalRing) error {
	for _, ring := range fractal.CooperationRings {
		for _, coinID := range ring.CoinIDs {
			if coin, ok := system.Coins[coinID]; !ok {
				return errors.New("coin not found")
			} else if coin.Status != models.Run {
				return errors.New("coin is not running")
			}
		}
	}
	return nil
}

func (system *System) informOthers(fractal *models.FractalRing) error {
	for _, ring := range fractal.CooperationRings {
		for _, coinID := range ring.CoinIDs {
			coin := system.Coins[coinID]
			coin.Status = models.Blocked
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

func (system *System) runFractal(fractal *models.FractalRing) error {
	for round := 0; round < models.RoundsCount; round++ {
		for index, ring := range fractal.CooperationRings {
			if ring.Rounds == -1 {
				totalErrors := []error{}
				for _, traderID := range fractal.VerificationTeam {
					if err := system.Traders[traderID].Vote(); err != nil {
						if err.Error() != "bad behavior" {
							return err
						}
						totalErrors = append(totalErrors, err)
					}
				}

				if 2*len(totalErrors) >= len(fractal.VerificationTeam) {
					ring.Rounds = round
					fractal.CooperationRings[index] = ring
					money := system.Coins[ring.CoinIDs[0]].Amount * float64(round) / models.RoundsCount
					if err := system.applyRing(ring, money); err != nil {
						return err
					}
				}
			}
		}
	}

	for index, ring := range fractal.CooperationRings {
		if ring.Rounds == -1 {
			ring.Rounds = models.RoundsCount
			fractal.CooperationRings[index] = ring
			if err := system.applyRing(ring, system.Coins[ring.CoinIDs[0]].Amount); err != nil {
				return err
			}
		}
	}
	return nil
}

func (system *System) applyRing(ring models.CooperationTable, money float64) error {
	for _, coinID := range ring.CoinIDs {
		coin := system.Coins[coinID]
		amount := money * coin.Amount / ring.Weight
		if ring.Rounds < models.RoundsCount {
			coin.Status = models.Expired
		} else {
			coin.Status = models.Paid
			amount += models.FractalPrize
		}
		system.Coins[coinID] = coin

		for _, trader := range system.Traders {
			if err := trader.UpdateBalance(coin.Owner, amount); err != nil {
				return err
			}
		}
	}

	for _, trader := range system.Traders {
		if ring.Rounds < models.RoundsCount {
			trader.ExpireRing(ring)
		} else {
			trader.PayRing(ring)
		}
	}
	return nil
}

func (system *System) CreateRandomCoins(trader *models.Trader, done <-chan bool, errors chan<- error) {
	for {
		select {
		case <-done:
			return
		case <-trader.Data.Ticker.C:
			amount := rand.Float64() * 10
			if trader.Account < amount {
				return
			}

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
			var trader *models.Trader
			amount, wallet := rand.Float64()*1000, uuid.New().String()
			if i < numRandomVoters {
				trader = models.CreateTrader(models.RandomVote, amount, wallet, coinTypeCount)
			} else if i < numRandomVoters+numBadVoters {
				trader = models.CreateTrader(models.BadVote, amount, wallet, coinTypeCount)
			} else {
				trader = models.CreateTrader(models.Normal, amount, wallet, coinTypeCount)
			}

			system.Locker.Lock()
			defer system.Locker.Unlock()
			system.Traders[trader.ID] = trader
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

	i, finished := 0, 0
	for _, trader := range system.Traders {
		dones[i] = make(chan bool, 1)
		go func(index int) {
			system.CreateRandomCoins(trader, dones[index], errors)
			finished++
		}(i)
		i++
	}

	for finished < len(system.Traders) {
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
				trader.Data.Ticker.Stop()
				dones[i] <- true
				i++
			}

			log.Println("Waiting for traders to finish...")
			time.Sleep(10 * time.Second)
			return
		}
	}
}

func (system *System) Save(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(system)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func Load(filePath string) (*System, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileStatus, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, fileStatus.Size())
	n, err := file.Read(data)
	if err != nil {
		return nil, err
	}

	system := NewSystem()
	if err := json.Unmarshal(data[:n], system); err != nil {
		return nil, err
	}
	return system, nil
}
