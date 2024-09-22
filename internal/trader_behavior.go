package internal

import (
	"fmt"
	"io"
	"log"
	"time"
)

const (
	Debug = false
)

func print(message string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), message)
}

func RandomBehavior(numTraders, numTypes, numCoins int) error {
	if !Debug {
		log.SetOutput(io.Discard)
	}
	print("Starting random behavior simulation...")

	traders, err := createTraders(numTraders, numTypes)
	if err != nil {
		return err
	}
	print("Traders created successfully.")

	coins, err := createCoins(numCoins, numTypes, traders)
	if err != nil {
		return err
	}
	print("Coins created successfully.")

	rings, fractals, err := processCoins(numCoins, traders, coins)
	if err != nil {
		return err
	}
	print("Coins processed successfully.")

	analyzeTraders(numTraders, rings, fractals)
	analyzeRings(numTypes, traders, rings)

	print("Random behavior simulation completed.")
	return nil
}
