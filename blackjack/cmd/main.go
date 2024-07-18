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

	// If dealer score <= 16, we hit
	// If dealer has a soft 17, then we hit.
	for blackjack.PlayerHand(dealer).Score() <= 16 || (blackjack.PlayerHand(dealer).Score() == 17 && blackjack.PlayerHand(dealer).MinScore() != 17) {
		pick, err := deck.Draw()
		if err != nil {
			panic(err)
		}

		dealer.Cards = append(dealer.Cards, pick)
	}

	pScore, dScore := player.Score(), blackjack.PlayerHand(dealer).Score()
	fmt.Println("Final hands")
	fmt.Printf("Player hand: %s, Score: %d\n", player, pScore)
	fmt.Printf("Dealer hand: %s, Score: %d\n", blackjack.PlayerHand(dealer), dScore)
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
