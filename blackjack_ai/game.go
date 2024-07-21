package blackjack_ai

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

	dealer   Hand
	dealerAI DealerAI
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
		deck:     deck_of_cards.NewDeck(deck_of_cards.Packs(opts.Decks), deck_of_cards.Shuffle),
		state:    statePlayerTurn,
		player:   Hand{cards: make([]deck_of_cards.Card, 0)},
		dealer:   Hand{cards: make([]deck_of_cards.Card, 0)},
		dealerAI: DealerAI{},
		opts:     opts,
	}
}

func (g *Game) Play(ai AI) int {
	for i := 1; i <= g.opts.Rounds; i++ {
		fmt.Printf("\nRound %d:\n", i)

		bet(g, ai)

		deal(g)

		for g.state == statePlayerTurn {
			move := ai.Play(g.player, g.dealer.cards[0])
			move(g)
		}

		for g.state == stateDealerTurn {
			move := g.dealerAI.Play(g.dealer, g.dealer.cards[0])
			move(g)
		}

		end(g, ai)
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

func bet(g *Game, ai AI) int {
	g.bet = ai.Bet()

	return g.bet
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

func end(g *Game, ai AI) {
	pScore, dScore := g.player.Score(), g.dealer.Score()
	winnings := g.bet

	var result string
	switch {
	case pScore > 21:
		winnings *= -1
		result = "You busted"
	case dScore > 21:
		winnings = int(float32(winnings) * g.opts.Payout)
		result = "Dealer busted"
	case pScore > dScore:
		winnings = int(float32(winnings) * g.opts.Payout)
		result = "You win!"
	case dScore > pScore:
		winnings *= -1
		result = "You lose"
	case dScore == pScore:
		winnings *= 0
		result = "Draw"
	}

	g.balance += winnings

	ai.Results(g.player, g.dealer, result)
}

func draw(hand *Hand, deck *deck_of_cards.Deck) *Hand {
	pick, err := deck.Draw()
	if err != nil {
		panic(err)
	}

	hand.cards = append(hand.cards, pick)

	return hand
}
