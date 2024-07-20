package deck_of_cards

import (
	"math/rand"
	"testing"
)

func TestNew(t *testing.T) {
	deck := NewDeck()
	// 13 ranks * 4 suits
	if len(deck.Cards) != 13*4 {
		t.Error("Wrong number of cards in a new deck.")
	}
}

func TestDefaultSort(t *testing.T) {
	deck := NewDeck(DefaultSort)
	exp := Card{Rank: Ace, Suit: Spade}
	if deck.Cards[0] != exp {
		t.Error("Expected Ace of Spades as first card. Received:", deck.Cards[0])
	}
}

func TestSort(t *testing.T) {
	deck := NewDeck(Sort(less))
	exp := Card{Rank: Ace, Suit: Spade}
	if deck.Cards[0] != exp {
		t.Error("Expected Ace of Spades as first card. Received:", deck.Cards[0])
	}
}

func TestShuffle(t *testing.T) {
	// make shuffleRand deterministic
	// First call to shuffleRand.Perm(52) should be:
	// [40 35 ... ]
	shuffleRand = rand.New(rand.NewSource(0))

	orig := NewDeck()
	first := orig.Cards[40]
	second := orig.Cards[35]

	shuffled := NewDeck(Shuffle)
	if shuffled.Cards[0] != first {
		t.Errorf("Expected the first card to be %s, received %s.", first, shuffled.Cards[0])
	}
	if shuffled.Cards[1] != second {
		t.Errorf("Expected the first card to be %s, received %s.", second, shuffled.Cards[1])
	}
}

func TestJokers(t *testing.T) {
	deck := NewDeck(Jokers(4))
	count := 0
	for _, c := range deck.Cards {
		if c.Suit == Joker {
			count++
		}
	}
	if count != 4 {
		t.Error("Expected 4 Jokers, received:", count)
	}
}

func TestFilter(t *testing.T) {
	filter := func(card Card) bool {
		return card.Rank == Two || card.Rank == Three
	}
	deck := NewDeck(Filter(filter))
	for _, c := range deck.Cards {
		if c.Rank == Two || c.Rank == Three {
			t.Error("Expected all twos and threes to be filtered out.")
		}
	}
}

func TestPacks(t *testing.T) {
	deck := NewDeck(Packs(3))
	// 13 ranks * 4 suits * 3 decks
	if len(deck.Cards) != 13*4*3 {
		t.Errorf("Expected %d cards, received %d cards.", 13*4*3, len(deck.Cards))
	}
}
