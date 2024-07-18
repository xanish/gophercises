package deck_of_cards

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type Deck struct {
	cards []Card
}

var shuffleRand = rand.New(rand.NewSource(time.Now().Unix()))

func NewDeck(opts ...func([]Card) []Card) Deck {
	var cards []Card
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}

	for _, opt := range opts {
		cards = opt(cards)
	}

	return Deck{cards: cards}
}

func (d *Deck) Draw() (Card, error) {
	if len(d.cards) == 0 {
		return Card{}, fmt.Errorf("no more cards to draw from")
	}

	top, rest := d.cards[0], d.cards[1:]
	d.cards = rest
	return top, nil
}

func (d *Deck) String() string {
	cards := strings.Builder{}
	for _, card := range d.cards {
		cards.WriteString(fmt.Sprintf("%s\n", card.String()))
	}

	return cards.String()
}

func Shuffle(cards []Card) []Card {
	ret := make([]Card, len(cards))

	perm := shuffleRand.Perm(len(cards))
	for i, j := range perm {
		ret[i] = cards[j]
	}

	return ret
}

func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, less(cards))
	return cards
}

func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

func Filter(f func(card Card) bool) func(cards []Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card

		for _, c := range cards {
			if !f(c) {
				ret = append(ret, c)
			}
		}
		return ret
	}
}

func Packs(count int) func(cards []Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for i := 0; i < count; i++ {
			ret = append(ret, cards...)
		}

		return ret
	}
}

func Jokers(count int) func(cards []Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < count; i++ {
			cards = append(cards, Card{Joker, Rank(i)})
		}

		return cards
	}
}

func less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return cards[i].AbsRank() < cards[j].AbsRank()
	}
}
