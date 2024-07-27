package blackjack_ai

import "github.com/xanish/gophercises/deck_of_cards"

type DealerAI struct{}

func (ai DealerAI) Bet(shuffled bool) int {
	return 1
}

func (ai DealerAI) Play(hand Hand, dealer deck_of_cards.Card) Move {
	// If dealer score <= 16, we hit
	// If dealer has a soft 17, then we hit.
	if hand.Score() <= 16 || (hand.Score() == 17 && hand.IsSoftScore()) {
		return Hit
	}

	return Stand
}

func (ai DealerAI) Results(hand Hand, dealer Hand) {
	// noop
}
