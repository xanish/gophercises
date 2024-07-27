package blackjack_ai

import (
	"fmt"
	"strings"

	"github.com/xanish/gophercises/deck_of_cards"
)

type Hand struct {
	cards []deck_of_cards.Card
	bet   int
}

func NewHand(cards []deck_of_cards.Card, bet int) Hand {
	return Hand{
		cards: cards,
		bet:   bet,
	}
}

func (h Hand) NumCards() int {
	return len(h.cards)
}

func (h Hand) Score() int {
	minScore := h.minScore()
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

func (h Hand) Scores() []int {
	var scores []int
	for _, c := range h.cards {
		scores = append(scores, NewHand([]deck_of_cards.Card{c}, 0).Score())
	}

	return scores
}

// IsSoftScore returns true if the score of a hand is a soft score - that is if an ace
// is being counted as 11 points.
func (h Hand) IsSoftScore() bool {
	return h.minScore() != h.Score()
}

func (h Hand) minScore() int {
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

	return fmt.Sprintf("[%s], (Score: %d)", strings.Join(cards, ", "), h.Score())
}
