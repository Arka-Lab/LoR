package main

import (
	"flag"
	"log"
	"time"

	"github.com/Arka-Lab/LoR/internal"
	"github.com/Arka-Lab/LoR/pkg"
)

func ParseFlags() (int, time.Duration, int, int, int, string) {
	typesPtr := flag.Int("type", 3, "number of coin types")
	runTimePtr := flag.Int("time", 60, "run time in seconds")
	tradersPtr := flag.Int("trader", 100, "number of traders")
	randomsPtr := flag.Int("random", 0, "number of random traders")
	badsPtr := flag.Int("bad", 0, "number of bad traders")
	filepathPtr := flag.String("file", "system.json", "file path to save system")
	flag.Parse()

	numTypes, numTraders := *typesPtr, *tradersPtr
	runTime := time.Duration(*runTimePtr) * time.Second
	numRandoms, numBads := *randomsPtr, *badsPtr
	filePath := *filepathPtr

	return numTypes, runTime, numTraders, numRandoms, numBads, filePath
}

func main() {
	logger := log.Default()
	finish := make(chan bool, 1)
	system := internal.NewSystem()

	numTypes, runTime, numTraders, numRandoms, numBads, filePath := ParseFlags()

	logger.Printf("Starting simulation with %d types (alpha = %.2f%%)...\n", numTypes, pkg.BadBehavior*100)
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

	if err := system.Save(filePath); err != nil {
		logger.Fatalf("Error saving system: %v\n", err)
	}
	logger.Printf("System saved to %s\n", filePath)

	internal.AnalyzeSystem(system)
}
