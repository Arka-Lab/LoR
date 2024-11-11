package main

import (
	"log"
	"time"

	"github.com/Arka-Lab/LoR/internal"
)

func main() {
	logger := log.Default()
	finish := make(chan bool, 1)
	system := internal.NewSystem()
	numTraders, numTypes, runTime := 100, 3, 90*time.Second

	logger.Printf("Starting system with %d traders and %d types...\n", numTraders, numTypes)
	system.Init(numTraders, uint(numTypes))
	logger.Println("System initialized!")

	logger.Println("Starting system...")
	done := make(chan bool, 1)
	go func() {
		system.Start(finish)
		done <- true
	}()
	logger.Println("System started!")

	logger.Printf("Waiting for %s...\n", runTime)
	time.Sleep(runTime)
	finish <- true
	<-done

	logger.Println("System stopped!")
	logger.Println("Number of coins:", len(system.Coins))
	logger.Println("Number of fractal rings:", len(system.Fractals))
}
