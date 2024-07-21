package main

import (
	"github.com/xanish/gophercises/blackjack_ai"
)

func main() {
	game := blackjack_ai.New(2)
	game.Play(blackjack_ai.NewPlayerAI())
}
