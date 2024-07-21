package main

import (
	"fmt"
	"github.com/xanish/gophercises/blackjack"
)

func main() {
	game := blackjack.New(blackjack.Options{Decks: 3, Rounds: 2})
	fmt.Printf("Your Balance: %d", game.Play())
}
