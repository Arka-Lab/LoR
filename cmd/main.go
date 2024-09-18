package main

import (
	"log"

	"github.com/Arka-Lab/LoR/internal"
)

func main() {
	numTraders, numTypes, numCoins := 100, 5, 200000
	if err := internal.RandomBehavior(numTraders, numTypes, numCoins); err != nil {
		log.Fatal(err)
	}
}
