package blackjack_ai

import (
	"fmt"
	"github.com/xanish/gophercises/deck_of_cards"
)

type PlayerAI struct{}

func NewPlayerAI() *PlayerAI {
	return &PlayerAI{}
}

func (ai *PlayerAI) Bet() int {
	return 1
}

func (ai *PlayerAI) Play(hand Hand, dealer deck_of_cards.Card) Move {
	for {
		var input string
		fmt.Println()
		fmt.Printf("Player's hand: %s\n", hand)
		fmt.Printf("Dealer's card: %s\n", dealer)
		fmt.Print("\nWhat would you like to do? (h)it or (s)tand... ")
		_, _ = fmt.Scanln(&input)

		switch input {
		case "h":
			return Hit
		case "s":
			return Stand
		default:
			fmt.Printf("Invalid choice: %s\n", input)
		}
	}
}

func (ai *PlayerAI) Results(hand Hand, dealer Hand, result string) {
	fmt.Println("\nFinal hands:")
	fmt.Printf("Player's hand: %s\n", hand)
	fmt.Printf("Dealer's hand: %s\n\n", dealer)
	fmt.Printf("Result: %s\n", result)
}
