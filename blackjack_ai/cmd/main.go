package main

import (
	"fmt"
	"github.com/xanish/gophercises/blackjack_ai"
	"github.com/xanish/gophercises/deck_of_cards"
)

type noobPlayerAI struct {
	score int
	seen  int
	decks int
}

func (ai *noobPlayerAI) Bet(shuffled bool) int {
	if shuffled {
		ai.score = 0
		ai.seen = 0
	}

	remainingDecks := (ai.decks*52 - ai.seen) / 52
	trueScore := ai.score / remainingDecks
	switch {
	case trueScore >= 14:
		return 100000
	case trueScore >= 8:
		return 5000
	default:
		return 100
	}
}

func (ai *noobPlayerAI) Play(hand blackjack_ai.Hand, dealer deck_of_cards.Card) blackjack_ai.Move {
	score := hand.Score()
	if hand.NumCards() == 2 {
		if blackjack_ai.CanSplit(hand) {
			cardScore := hand.Score() / 2

			// card is ace, 8 or 9
			if cardScore >= 8 && cardScore != 10 {
				return blackjack_ai.Split
			}
		}

		if (score == 10 || score == 11) && !blackjack_ai.IsSoft(hand) {
			return blackjack_ai.Double
		}
	}

	dScore := blackjack_ai.NewHand([]deck_of_cards.Card{dealer}, 0).Score()
	if dScore >= 5 && dScore <= 6 {
		return blackjack_ai.Stand
	}

	if score < 13 {
		return blackjack_ai.Hit
	}

	return blackjack_ai.Stand
}

func (ai *noobPlayerAI) Results(hand blackjack_ai.Hand, dealer blackjack_ai.Hand, result string) {
	for _, score := range dealer.Scores() {
		ai.count(score)
	}

	for _, score := range hand.Scores() {
		ai.count(score)
	}
}

func (ai *noobPlayerAI) count(score int) {
	switch {
	case score >= 10:
		ai.score--
	case score <= 6:
		ai.score++
	}
	ai.seen++
}

func main() {
	game := blackjack_ai.New(blackjack_ai.Options{Decks: 3})
	fmt.Printf("Your Balance: %d\n", game.Play(&noobPlayerAI{
		decks: 3,
	}))
}
