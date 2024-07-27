package blackjack_ai

import (
	"github.com/xanish/gophercises/deck_of_cards"
)

type AI interface {
	Bet(shuffled bool) int
	Play(hand Hand, dealer deck_of_cards.Card) Move
	Results(hand []Hand, dealer Hand)
}

type Move func(*Game) error
