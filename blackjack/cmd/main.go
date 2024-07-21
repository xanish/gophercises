package main

import (
	"fmt"
	"github.com/xanish/gophercises/blackjack"
)

func main() {
	g := blackjack.New()

	for round := 1; round <= 10; round++ {
		fmt.Printf("Round %d:\n", round)

		blackjack.Deal(&g)

		var input string
		for g.State == blackjack.StatePlayerTurn {
			fmt.Println()
			fmt.Printf("Player hand: %s\n", g.Player)
			fmt.Printf("Dealer hand: %s\n", g.Dealer.DealerString())
			fmt.Print("\nWhat would you like to do? (h)it or (s)tand... ")
			_, _ = fmt.Scanln(&input)

			switch input {
			case "h":
				blackjack.Hit(&g)
			case "s":
				blackjack.Stand(&g)
			default:
				fmt.Printf("Invalid choice: %s\n", input)
			}
		}

		for g.State == blackjack.StatePlayerTurn {
			// If dealer score <= 16, we hit
			// If dealer has a soft 17, then we hit.
			if g.Dealer.Score() <= 16 || (g.Dealer.Score() == 17 && g.Dealer.MinScore() != 17) {
				blackjack.Hit(&g)
			} else {
				blackjack.Stand(&g)
			}
		}

		blackjack.End(&g)
	}
}
