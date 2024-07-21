package blackjack

import (
	"fmt"
	"github.com/xanish/gophercises/deck_of_cards"
)

type state int8

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

type Game struct {
	deck   deck_of_cards.Deck
	state  state
	player Hand
	dealer Hand
}

func New() Game {
	return Game{
		deck:   deck_of_cards.NewDeck(deck_of_cards.Packs(3), deck_of_cards.Shuffle),
		player: Hand{cards: make([]deck_of_cards.Card, 0)},
		dealer: Hand{cards: make([]deck_of_cards.Card, 0)},
		state:  statePlayerTurn,
	}
}

func (g *Game) Play() {
	deal(g)

	var input string
	for g.state == statePlayerTurn {
		fmt.Println()
		fmt.Printf("player hand: %s\n", g.player)
		fmt.Printf("dealer hand: %s\n", g.dealer.DealerString())
		fmt.Print("\nWhat would you like to do? (h)it or (s)tand... ")
		_, _ = fmt.Scanln(&input)

		switch input {
		case "h":
			Hit(g)
		case "s":
			Stand(g)
		default:
			fmt.Printf("Invalid choice: %s\n", input)
		}
	}

	for g.state == statePlayerTurn {
		// If dealer score <= 16, we hit
		// If dealer has a soft 17, then we hit.
		if g.dealer.Score() <= 16 || (g.dealer.Score() == 17 && g.dealer.MinScore() != 17) {
			Hit(g)
		} else {
			Stand(g)
		}
	}

	end(g)
}

func (g *Game) String() string {
	return fmt.Sprintf("current deck size: %d\n state: %v\n player hand: %s\n dealer hand: %s\n", g.deck.RemainingCards(), g.state, g.player, g.dealer)
}

func (g *Game) currentPlayer() *Hand {
	switch g.state {
	case statePlayerTurn:
		return &g.player
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("it isn't currently any player's turn")
	}
}

func Hit(g *Game) {
	pick, err := g.deck.Draw()
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
	g.state++
}

func deal(g *Game) {
	g.player.cards = make([]deck_of_cards.Card, 0, 2)
	g.dealer.cards = make([]deck_of_cards.Card, 0, 2)

	for i := 0; i < 2; i++ {
		draw(&g.player, &g.deck)
		draw(&g.dealer, &g.deck)
	}

	g.state = statePlayerTurn
}

func end(g *Game) {
	pScore, dScore := g.player.Score(), g.dealer.Score()
	fmt.Println("\nFinal hands:")
	fmt.Printf("player hand: %s (Score: %d)\n", g.player, pScore)
	fmt.Printf("dealer hand: %s (Score: %d)\n\n", g.dealer, dScore)
	switch {
	case pScore > 21:
		fmt.Println("Result: You busted")
	case dScore > 21:
		fmt.Println("Result: dealer busted")
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
