package main

import (
	"fmt"
	"github.com/xanish/gophercises/blackjack_ai"
)

func main() {
	game := blackjack_ai.New(blackjack_ai.Options{Decks: 3, Rounds: 2})
	fmt.Printf("Your Balance: %d", game.Play(blackjack_ai.NewPlayerAI()))
}
