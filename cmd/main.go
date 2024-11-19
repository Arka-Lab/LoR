package main

import (
	"log"
	"time"

	"github.com/Arka-Lab/LoR/internal"
	"github.com/Arka-Lab/LoR/pkg"
)

// TODO: Note that the number of traders are independent of the results and we used small numbers cause of resource limitations.

func main() {
	logger := log.Default()
	finish := make(chan bool, 1)
	system := internal.NewSystem()
	numTypes, runTime := 3, 90*time.Second
	numTraders, numRandoms, numBads := 100, 20, 20

	logger.Printf("Starting simulation with %d types (BadBehavior percentage = %.2f%%)...\n", numTypes, pkg.BadBehavior*100)
	system.Init(numTraders, numRandoms, numBads, uint(numTypes))
	logger.Println("Simulation initialized!")

	logger.Println("Starting simulation...")
	done := make(chan bool, 1)
	go func() {
		system.Start(finish)
		done <- true
	}()
	logger.Println("Simulation started!")

	logger.Printf("Waiting for %s...\n", runTime)
	time.Sleep(runTime)
	finish <- true
	<-done
	logger.Println("Simulation stopped!")

	internal.AnalyzeSystem(system)
}
