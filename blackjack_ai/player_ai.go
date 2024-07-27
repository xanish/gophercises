package blackjack_ai

import (
	"fmt"
	"github.com/xanish/gophercises/deck_of_cards"
)

type PlayerAI struct{}

func NewPlayerAI() *PlayerAI {
	return &PlayerAI{}
}

func (ai *PlayerAI) Bet(shuffled bool) int {
	if shuffled {
		fmt.Println("The deck was just shuffled.")
	}

	fmt.Print("How much will you bet? ")
	var bet int
	_, _ = fmt.Scanln(&bet)

	return bet
}

func (ai *PlayerAI) Play(hand Hand, dealer deck_of_cards.Card) Move {
	for {
		var input string
		fmt.Println()
		fmt.Printf("Player's hand: %s\n", hand)
		fmt.Printf("Dealer's card: %s\n", dealer)
		fmt.Print("\nWhat would you like to do? (h)it, (s)tand, (d)ouble or s(p)lit... ")
		_, _ = fmt.Scanln(&input)

		switch input {
		case "h":
			return Hit
		case "s":
			return Stand
		case "d":
			return Double
		case "p":
			return Split
		default:
			fmt.Printf("Invalid choice: %s\n", input)
		}
	}
}

func (ai *PlayerAI) Results(hand Hand, dealer Hand, result string) {
	fmt.Printf("Player's hand: %s\n", hand)
	fmt.Printf("Dealer's hand: %s\n", dealer)
	fmt.Printf("Result: %s\n\n", result)
}
