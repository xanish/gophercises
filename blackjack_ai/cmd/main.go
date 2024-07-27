package main

import (
	"fmt"
	"github.com/xanish/gophercises/blackjack_ai"
)

func main() {
	game := blackjack_ai.New(blackjack_ai.Options{Decks: 1})
	fmt.Printf("Your Balance: %d\n", game.Play(blackjack_ai.NewPlayerAI()))
}
