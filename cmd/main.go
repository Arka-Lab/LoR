package main

import (
	"log"
	"time"

	"github.com/Arka-Lab/LoR/internal"
	"github.com/Arka-Lab/LoR/pkg"
)

func main() {
	logger := log.Default()
	finish := make(chan bool, 1)
	system := internal.NewSystem()
	numTypes, runTime := 3, 90*time.Second
	numTraders, numRandoms, numBads := 100, 10, 10

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
