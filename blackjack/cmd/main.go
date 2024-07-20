package main

import (
	"fmt"
	"github.com/xanish/gophercises/blackjack"
)

func main() {
	gs := blackjack.NewGameState()

	for round := 1; round <= 10; round++ {
		fmt.Printf("Round %d:\n", round)

		gs = blackjack.Deal(gs)

		var input string
		for gs.State == blackjack.StatePlayerTurn {
			fmt.Println()
			fmt.Printf("Player hand: %s\n", gs.Player)
			fmt.Printf("Dealer hand: %s\n", gs.Dealer.DealerString())
			fmt.Print("\nWhat would you like to do? (h)it or (s)tand... ")
			_, _ = fmt.Scanln(&input)

			switch input {
			case "h":
				gs = blackjack.Hit(gs)
			case "s":
				gs = blackjack.Stand(gs)
			default:
				fmt.Printf("Invalid choice: %s\n", input)
			}
		}

		for gs.State == blackjack.StatePlayerTurn {
			// If dealer score <= 16, we hit
			// If dealer has a soft 17, then we hit.
			if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
				gs = blackjack.Hit(gs)
			} else {
				gs = blackjack.Stand(gs)
			}
		}

		gs = blackjack.End(gs)
	}
}
