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

func New() Game {
	return Game{
		Deck:   deck_of_cards.NewDeck(deck_of_cards.Packs(3), deck_of_cards.Shuffle),
		Player: Hand{cards: make([]deck_of_cards.Card, 0)},
		Dealer: Hand{cards: make([]deck_of_cards.Card, 0)},
		State:  StatePlayerTurn,
	}
}

func (g *Game) String() string {
	return fmt.Sprintf("current deck size: %d\n state: %v\n player hand: %s\n dealer hand: %s\n", g.Deck.RemainingCards(), g.State, g.Player, g.Dealer)
}

func (g *Game) currentPlayer() *Hand {
	switch g.State {
	case StatePlayerTurn:
		return &g.Player
	case StateDealerTurn:
		return &g.Dealer
	default:
		panic("it isn't currently any player's turn")
	}
}

func Deal(g *Game) {
	g.Player.cards = make([]deck_of_cards.Card, 0, 2)
	g.Dealer.cards = make([]deck_of_cards.Card, 0, 2)

	for i := 0; i < 2; i++ {
		draw(&g.Player, &g.Deck)
		draw(&g.Dealer, &g.Deck)
	}

	g.State = StatePlayerTurn
}

func Hit(g *Game) {
	pick, err := g.Deck.Draw()
	if err != nil {
		panic(err)
	}

	player := g.currentPlayer()
	player.cards = append(player.cards, pick)

	if player.Score() > 21 {
		Stand(g)
	}
}

func Stand(g *Game) {
	g.State++
}

func End(g *Game) {
	pScore, dScore := g.Player.Score(), g.Dealer.Score()
	fmt.Println("\nFinal hands:")
	fmt.Printf("Player hand: %s (Score: %d)\n", g.Player, pScore)
	fmt.Printf("Dealer hand: %s (Score: %d)\n\n", g.Dealer, dScore)
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
		fmt.Println("Result: draw")
	}
}

func draw(hand *Hand, deck *deck_of_cards.Deck) *Hand {
	pick, err := deck.Draw()
	if err != nil {
		panic(err)
	}

	hand.cards = append(hand.cards, pick)

	return hand
}
