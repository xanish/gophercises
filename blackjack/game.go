package blackjack

import (
	"fmt"
	"github.com/xanish/gophercises/deck_of_cards"
)

type State int8

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

type Game struct {
	Deck   deck_of_cards.Deck
	State  State
	Player Hand
	Dealer Hand
}

func (g *Game) String() string {
	return fmt.Sprintf("current deck size: %d\n state: %v\n player hand: %s\n dealer hand: %s\n", g.Deck.RemainingCards(), g.State, g.Player, g.Dealer)
}

func (g *Game) CurrentPlayer() *Hand {
	switch g.State {
	case StatePlayerTurn:
		return &g.Player
	case StateDealerTurn:
		return &g.Dealer
	default:
		panic("it isn't currently any player's turn")
	}
}

func NewGameState() Game {
	return Game{
		Deck:   deck_of_cards.NewDeck(deck_of_cards.Packs(3), deck_of_cards.Shuffle),
		Player: Hand{Cards: make([]deck_of_cards.Card, 0)},
		Dealer: Hand{Cards: make([]deck_of_cards.Card, 0)},
		State:  StatePlayerTurn,
	}
}

func Draw(hand *Hand, deck *deck_of_cards.Deck) *Hand {
	pick, err := deck.Draw()
	if err != nil {
		panic(err)
	}

	hand.Cards = append(hand.Cards, pick)

	return hand
}

func Deal(g Game) Game {
	ret := clone(g)

	ret.Player.Cards = make([]deck_of_cards.Card, 0, 2)
	ret.Dealer.Cards = make([]deck_of_cards.Card, 0, 2)

	for i := 0; i < 2; i++ {
		Draw(&ret.Player, &ret.Deck)
		Draw(&ret.Dealer, &ret.Deck)
	}

	ret.State = StatePlayerTurn

	return ret
}

func Hit(g Game) Game {
	ret := clone(g)

	pick, err := ret.Deck.Draw()
	if err != nil {
		panic(err)
	}

	player := ret.CurrentPlayer()
	player.Cards = append(player.Cards, pick)

	if player.Score() > 21 {
		return Stand(ret)
	}

	return ret
}

func Stand(g Game) Game {
	ret := clone(g)
	ret.State++

	return ret
}

func End(g Game) Game {
	ret := clone(g)

	pScore, dScore := ret.Player.Score(), ret.Dealer.Score()
	fmt.Println("\nFinal hands:")
	fmt.Printf("Player hand: %s (Score: %d)\n", ret.Player, pScore)
	fmt.Printf("Dealer hand: %s (Score: %d)\n\n", ret.Dealer, dScore)
	switch {
	case pScore > 21:
		fmt.Println("Result: You busted")
	case dScore > 21:
		fmt.Println("Result: Dealer busted")
	case pScore > dScore:
		fmt.Println("Result: You win!")
	case dScore > pScore:
		fmt.Println("Result: You lose")
	case dScore == pScore:
		fmt.Println("Result: Draw")
	}

	return ret
}

func clone(g Game) Game {
	ret := Game{
		Deck:   deck_of_cards.From(g.Deck),
		Player: Hand{Cards: make([]deck_of_cards.Card, len(g.Player.Cards))},
		Dealer: Hand{Cards: make([]deck_of_cards.Card, len(g.Dealer.Cards))},
		State:  g.State,
	}

	copy(ret.Player.Cards, g.Player.Cards)
	copy(ret.Dealer.Cards, g.Dealer.Cards)

	return ret
}
