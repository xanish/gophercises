package main

import (
	"fmt"
	"github.com/xanish/gophercises/blackjack"
	"github.com/xanish/gophercises/deck_of_cards"
)

func main() {
	deck := deck_of_cards.NewDeck(deck_of_cards.Packs(3), deck_of_cards.Shuffle)
	var player blackjack.PlayerHand
	var dealer blackjack.DealerHand

	for i := 0; i < 2; i++ {
		// distribute to players
		for _, hand := range []*blackjack.PlayerHand{&player} {
			pick, err := deck.Draw()
			if err != nil {
				panic(err)
			}

			hand.Cards = append(hand.Cards, pick)
		}

		// distribute to dealer
		pick, err := deck.Draw()
		if err != nil {
			panic(err)
		}

		dealer.Cards = append(dealer.Cards, pick)
	}

	var input string
	for input != "s" {
		fmt.Println(player)
		fmt.Println(dealer)
		fmt.Println("What would you like to do? (h)it or (s)tand...")
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

	fmt.Println("Final hands")
	fmt.Println(player)
	fmt.Println(blackjack.PlayerHand(dealer))
}
