package blackjack_ai

import (
	"errors"
	"fmt"
	"github.com/xanish/gophercises/deck_of_cards"
)

type state int8

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

var (
	errBust = errors.New("hand score exceeded 21")
)

type Game struct {
	deck  deck_of_cards.Deck
	state state
	opts  Options

	player               []Hand
	currentHandIdxInPlay int
	balance              int

	dealer   Hand
	dealerAI DealerAI
}

type Options struct {
	Decks  int
	Payout float32
}

func New(opts Options) Game {
	if opts.Decks == 0 {
		opts.Decks = 3
	}

	if opts.Payout == 0.0 {
		opts.Payout = 1.5
	}

	hand := make([]Hand, 0)
	hand = append(hand, Hand{cards: make([]deck_of_cards.Card, 0)})
	return Game{
		deck:     deck_of_cards.NewDeck(deck_of_cards.Packs(opts.Decks), deck_of_cards.Shuffle),
		state:    statePlayerTurn,
		player:   hand,
		dealer:   Hand{cards: make([]deck_of_cards.Card, 0)},
		dealerAI: DealerAI{},
		opts:     opts,
	}
}

func (g *Game) Play(ai AI) int {
	bet(g, ai, false)

	deal(g)

	if isBlackJack(g.dealer) {
		end(g, ai)
	}

	for i := 0; i < len(g.player); i++ {
		g.currentHandIdxInPlay = i

		if float32(g.deck.RemainingCards()) < float32(g.opts.Decks*52)*0.5 {
			g.deck = deck_of_cards.NewDeck(deck_of_cards.Packs(g.opts.Decks), deck_of_cards.Shuffle)
			bet(g, ai, true)
		}

		for g.state == statePlayerTurn && g.currentHandIdxInPlay == i {
			move := ai.Play(g.player[i], g.dealer.cards[0])
			err := move(g)

			bust := false
			for err != nil {
				// go to the next split
				if errors.Is(err, errBust) {
					bust = true
					break
				}

				fmt.Printf("\nInvalid move: %v\n", err)
				move = ai.Play(g.player[i], g.dealer.cards[0])
				err = move(g)
			}

			if bust {
				fmt.Println("Busted")
				break
			}
		}
	}

	for g.state == stateDealerTurn {
		move := g.dealerAI.Play(g.dealer, g.dealer.cards[0])
		_ = move(g)
	}

	end(g, ai)

	return g.balance
}

func (g *Game) String() string {
	return fmt.Sprintf("Remaining Cards In Deck: %d\n state: %v\n Player's hand: %s\n Dealer's hand: %s\n", g.deck.RemainingCards(), g.state, g.player, g.dealer)
}

func (g *Game) currentHandInPlay() *Hand {
	switch g.state {
	case statePlayerTurn:
		return &g.player[g.currentHandIdxInPlay]
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("it isn't currently any player's turn")
	}
}

func Hit(g *Game) error {
	pick, err := g.deck.Draw()
	if err != nil {
		panic(err)
	}

	player := g.currentHandInPlay()
	player.cards = append(player.cards, pick)

	if player.Score() > 21 {
		return errBust
	}

	return nil
}

func Stand(g *Game) error {
	if g.state == stateDealerTurn {
		g.state++
		return nil
	}

	if g.state == statePlayerTurn {
		g.currentHandIdxInPlay++
		if g.currentHandIdxInPlay >= len(g.player) {
			g.state++
		}
		return nil
	}

	return errors.New("invalid game state")
}

func Split(g *Game) error {
	hand := g.currentHandInPlay()
	if len(hand.cards) != 2 {
		return errors.New("you can only split with two cards in your hand")
	}

	if hand.cards[0].Rank != hand.cards[1].Rank {
		return errors.New("both cards must have the same rank to split")
	}

	// make a new hand out of the second hard in player's hand
	g.player = append(g.player, Hand{
		cards: []deck_of_cards.Card{hand.cards[1]},
		bet:   hand.bet,
	})

	// make the old hand have only the first card since we split the
	// second one in new hand
	g.player[g.currentHandIdxInPlay].cards = []deck_of_cards.Card{hand.cards[0]}

	return nil
}

func Double(g *Game) error {
	hand := g.currentHandInPlay()
	if len(hand.cards) != 2 {
		return errors.New("can only double on a hand with 2 cards")
	}

	hand.bet *= 2

	_ = Hit(g)
	_ = Stand(g)
	return nil
}

func bet(g *Game, ai AI, shuffled bool) int {
	hand := g.currentHandInPlay()
	hand.bet = ai.Bet(shuffled)

	return hand.bet
}

func deal(g *Game) {
	hand := g.currentHandInPlay()
	hand.cards = make([]deck_of_cards.Card, 0, 2)
	g.dealer.cards = make([]deck_of_cards.Card, 0, 2)

	for i := 0; i < 2; i++ {
		draw(hand, &g.deck)
		draw(&g.dealer, &g.deck)
	}

	g.state = statePlayerTurn
}

func end(g *Game, ai AI) {
	fmt.Println("\nFinal hands\n")

	for i, hand := range g.player {
		pScore, pBlackJack := hand.Score(), isBlackJack(hand)
		dScore, dBlackJack := g.dealer.Score(), isBlackJack(g.dealer)
		winnings := hand.bet

		var result string
		switch {
		case pBlackJack && dBlackJack:
			winnings = 0
			result = "Draw"
		case dBlackJack:
			winnings = -winnings
			result = "Dealer BlackJack"
		case pBlackJack:
			winnings = int(float32(winnings) * g.opts.Payout)
			result = "BlackJack"
		case pScore > 21:
			winnings = -winnings
			result = "Busted"
		case dScore > 21:
			result = "Dealer Busted"
		case pScore > dScore:
			result = "Won"
		case dScore > pScore:
			winnings = -winnings
			result = "Lost"
		case dScore == pScore:
			winnings = 0
			result = "Draw"
		}

		if len(g.player) > 1 {
			fmt.Printf("Hand %d:\n", i+1)
		}
		ai.Results(hand, g.dealer, result)

		g.balance += winnings
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

func isBlackJack(hand Hand) bool {
	return len(hand.cards) == 2 && hand.Score() == 21
}
