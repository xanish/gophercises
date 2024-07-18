package main

import (
	"fmt"

	"github.com/xanish/gophercises/blackjack"
	"github.com/xanish/gophercises/deck_of_cards"
)

func main() {
	deck := deck_of_cards.NewDeck(deck_of_cards.Packs(3), deck_of_cards.Shuffle)
	var player blackjack.Hand
	var dealer blackjack.Hand

	for i := 0; i < 2; i++ {
		for _, hand := range []*blackjack.Hand{&player, &dealer} {
			pick, err := deck.Draw()
			if err != nil {
				panic(err)
			}

			hand.Cards = append(hand.Cards, pick)
		}
	}

	var input string
	for input != "s" {
		fmt.Printf("Player hand: %s\n", player)
		fmt.Printf("Dealer hand: %s\n", dealer)
		fmt.Print("\nWhat would you like to do? (h)it or (s)tand... ")
		_, _ = fmt.Scanln(&input)

		switch input {
		case "h":
			pick, err := deck.Draw()
			if err != nil {
				panic(err)
			}

			player.Cards = append(player.Cards, pick)
		}
	}

	// If dealer score <= 16, we hit
	// If dealer has a soft 17, then we hit.
	for dealer.Score() <= 16 || (dealer.Score() == 17 && dealer.MinScore() != 17) {
		pick, err := deck.Draw()
		if err != nil {
			panic(err)
		}

		dealer.Cards = append(dealer.Cards, pick)
	}

	pScore, dScore := player.Score(), dealer.Score()
	fmt.Println("\nFinal hands")
	fmt.Printf("Player hand: %s, Score: %d\n", player, pScore)
	fmt.Printf("Dealer hand: %s, Score: %d\n", dealer, dScore)
	switch {
	case pScore > 21:
		fmt.Println("You busted")
	case dScore > 21:
		fmt.Println("Dealer busted")
	case pScore > dScore:
		fmt.Println("You win!")
	case dScore > pScore:
		fmt.Println("You lose")
	case dScore == pScore:
		fmt.Println("Draw")
	}
}
