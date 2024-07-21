package blackjack

import (
	"strings"

	"github.com/xanish/gophercises/deck_of_cards"
)

type Hand struct {
	cards []deck_of_cards.Card
}

func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}

	for _, c := range h.cards {
		if c.Rank == deck_of_cards.Ace {
			// ace is currently worth 1, and we are changing it to be worth 11
			// 11 - 1 = 10
			return minScore + 10
		}
	}

	return minScore
}

func (h Hand) MinScore() int {
	score := 0
	for _, c := range h.cards {
		score += min(int(c.Rank), 10)
	}
	return score
}

func (h Hand) String() string {
	var cards []string
	for _, card := range h.cards {
		cards = append(cards, card.String())
	}

	return strings.Join(cards, ", ")
}

func (h Hand) DealerString() string {
	var cards []string
	for i, card := range h.cards {
		if i > 0 {
			cards = append(cards, "**Hidden**")
		} else {
			cards = append(cards, card.String())
		}
	}

	return strings.Join(cards, ", ")
}
