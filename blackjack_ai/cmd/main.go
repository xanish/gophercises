package main

import (
	"github.com/xanish/gophercises/blackjack_ai"
)

func main() {
	game := blackjack_ai.New(blackjack_ai.Options{Decks: 3, Rounds: 2})
	game.Play(blackjack_ai.NewPlayerAI())
}
