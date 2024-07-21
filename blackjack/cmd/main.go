package main

import (
	"github.com/xanish/gophercises/blackjack"
)

func main() {
	g := blackjack.New(blackjack.Options{Decks: 3, Rounds: 2})

	g.Play()
}
