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
	deck  deck_of_cards.Deck
	state state
	opts  Options

	player  Hand
	bet     int
	balance int

	dealer Hand
}

type Options struct {
	Decks  int
	Rounds int
	Payout float32
}

func New(opts Options) Game {
	if opts.Decks == 0 {
		opts.Decks = 3
	}

	if opts.Rounds == 0 {
		opts.Rounds = 100
	}

	if opts.Payout == 0.0 {
		opts.Payout = 1.5
	}

	return Game{
		deck:   deck_of_cards.NewDeck(deck_of_cards.Packs(opts.Decks), deck_of_cards.Shuffle),
		player: Hand{cards: make([]deck_of_cards.Card, 0)},
		dealer: Hand{cards: make([]deck_of_cards.Card, 0)},
		state:  statePlayerTurn,
		opts:   opts,
	}
}

func (g *Game) Play() int {
	for i := 1; i <= g.opts.Rounds; i++ {
		fmt.Printf("Round %d:\n", i)

		shuffled := false
		if float32(g.deck.RemainingCards()) < float32(g.opts.Decks*52)*0.5 {
			g.deck = deck_of_cards.NewDeck(deck_of_cards.Packs(g.opts.Decks), deck_of_cards.Shuffle)
			shuffled = false
		}

		bet(g, shuffled)

		deal(g)

		var input string
		for g.state == statePlayerTurn {
			fmt.Println()
			fmt.Printf("Player's hand: %s\n", g.player)
			fmt.Printf("Dealer's hand: %s\n", g.dealer.DealerString())
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

		for g.state == stateDealerTurn {
			// If dealer score <= 16, we hit
			// If dealer has a soft 17, then we hit.
			if g.dealer.Score() <= 16 || (g.dealer.Score() == 17 && g.dealer.IsSoftScore()) {
				Hit(g)
			} else {
				Stand(g)
			}
		}

		end(g)
	}

	return g.balance
}

func (g *Game) String() string {
	return fmt.Sprintf("Remaining Cards In Deck: %d\n state: %v\n Player's hand: %s\n Dealer's hand: %s\n", g.deck.RemainingCards(), g.state, g.player, g.dealer)
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

func bet(g *Game, shuffled bool) {
	if shuffled {
		fmt.Println("The deck was just shuffled.")
	}

	fmt.Print("How much will you bet? ")
	_, _ = fmt.Scanln(&g.bet)
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
	winnings := g.bet

	fmt.Println("\nFinal hands:")
	fmt.Printf("Player's hand: %s (Score: %d)\n", g.player, pScore)
	fmt.Printf("Dealer's hand: %s (Score: %d)\n\n", g.dealer, dScore)
	switch {
	case pScore > 21:
		winnings *= -1
		fmt.Println("Result: You busted")
	case dScore > 21:
		winnings = int(float32(winnings) * g.opts.Payout)
		fmt.Println("Result: Dealer busted")
	case pScore > dScore:
		winnings = int(float32(winnings) * g.opts.Payout)
		fmt.Println("Result: You win!")
	case dScore > pScore:
		winnings *= -1
		fmt.Println("Result: You lose")
	case dScore == pScore:
		winnings *= 0
		fmt.Println("Result: Draw")
	}

	g.balance += winnings
}

func draw(hand *Hand, deck *deck_of_cards.Deck) *Hand {
	pick, err := deck.Draw()
	if err != nil {
		panic(err)
	}

	hand.cards = append(hand.cards, pick)

	return hand
}
