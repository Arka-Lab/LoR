package main

import (
	"flag"
	"log"
	"time"

	"github.com/Arka-Lab/LoR/internal"
	"github.com/Arka-Lab/LoR/pkg"
)

func ParseFlags() (int, time.Duration, int, int, int) {
	typesPtr := flag.Int("type", 3, "number of coin types")
	runTimePtr := flag.Int("time", 90, "run time in seconds")
	tradersPtr := flag.Int("trader", 100, "number of traders")
	randomsPtr := flag.Int("random", 30, "number of random traders")
	badsPtr := flag.Int("bad", 30, "number of bad traders")
	flag.Parse()

	numTypes, numTraders := *typesPtr, *tradersPtr
	runTime := time.Duration(*runTimePtr) * time.Second
	numRandoms, numBads := *randomsPtr, *badsPtr

	return numTypes, runTime, numTraders, numRandoms, numBads
}

func main() {
	logger := log.Default()
	finish := make(chan bool, 1)
	system := internal.NewSystem()

	numTypes, runTime, numTraders, numRandoms, numBads := ParseFlags()

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
