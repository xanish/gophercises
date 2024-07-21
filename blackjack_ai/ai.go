package blackjack_ai

import (
	"github.com/xanish/gophercises/deck_of_cards"
)

type AI interface {
	Bet() int
	Play(hand Hand, dealer deck_of_cards.Card) Move
	Results(hand Hand, dealer Hand, result string)
}

type Move func(*Game)
