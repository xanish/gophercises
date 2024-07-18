package blackjack

import (
	"github.com/xanish/gophercises/deck_of_cards"
	"strings"
)

type PlayerHand struct {
	Cards []deck_of_cards.Card
}

type DealerHand PlayerHand

func (h PlayerHand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}

	for _, c := range h.Cards {
		if c.Rank == deck_of_cards.Ace {
			// ace is currently worth 1, and we are changing it to be worth 11
			// 11 - 1 = 10
			return minScore + 10
		}
	}

	return minScore
}

func (h PlayerHand) MinScore() int {
	score := 0
	for _, c := range h.Cards {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
