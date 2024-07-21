package main

import (
	"fmt"
	"github.com/xanish/gophercises/blackjack"
)

func main() {
	g := blackjack.New()

	for round := 1; round <= 10; round++ {
		fmt.Printf("Round %d:\n", round)

		g.Play()
	}
}
