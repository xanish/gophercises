package blackjack

import (
	"github.com/xanish/gophercises/deck_of_cards"
	"strings"
)

type PlayerHand struct {
	Cards []deck_of_cards.Card
}

type DealerHand PlayerHand

func (h PlayerHand) String() string {
	var cards []string
	for _, card := range h.Cards {
		cards = append(cards, card.String())
	}

	return strings.Join(cards, ", ")
}

func (h DealerHand) String() string {
	var cards []string
	for i, card := range h.Cards {
		if i > 0 {
			cards = append(cards, "Hidden")
		} else {
			cards = append(cards, card.String())
		}
	}

	return strings.Join(cards, ", ")
}
