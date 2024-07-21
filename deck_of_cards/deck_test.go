package deck_of_cards

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	deck := NewDeck()
	// 13 ranks * 4 suits
	if len(deck.cards) != 13*4 {
		t.Error("Wrong number of cards in a new deck.")
	}
}

func TestFrom(t *testing.T) {
	deck := NewDeck()
	for i := 0; i < 5; i++ {
		_, err := deck.Draw()
		if err != nil {
			t.Fatalf("expected err to be nil, got %v", err)
		}
	}
	gotDeck := From(deck)
	want := deck.RemainingCards()
	got := gotDeck.RemainingCards()
	if got != want {
		t.Error("expected remaining cards to be the same.")
	}

	if !reflect.DeepEqual(deck, gotDeck) {
		t.Errorf("got deck %v; want %v", gotDeck, deck)
	}
}

func TestDeck_RemainingCards(t *testing.T) {
	deck := NewDeck()
	length := deck.RemainingCards()
	for i := 0; i < 5; i++ {
		_, err := deck.Draw()
		if err != nil {
			t.Fatalf("expected err to be nil, got %v", err)
		}
	}

	got := deck.RemainingCards()
	want := length - 5

	if got != want {
		t.Errorf("expected remaining cards to be %d, got %d", want, got)
	}
}

func TestDefaultSort(t *testing.T) {
	deck := NewDeck(DefaultSort)
	exp := Card{Rank: Ace, Suit: Spade}
	if deck.cards[0] != exp {
		t.Error("Expected Ace of Spades as first card. Received:", deck.cards[0])
	}
}

func TestSort(t *testing.T) {
	deck := NewDeck(Sort(less))
	exp := Card{Rank: Ace, Suit: Spade}
	if deck.cards[0] != exp {
		t.Error("Expected Ace of Spades as first card. Received:", deck.cards[0])
	}
}

func TestShuffle(t *testing.T) {
	// make shuffleRand deterministic
	// First call to shuffleRand.Perm(52) should be:
	// [40 35 ... ]
	shuffleRand = rand.New(rand.NewSource(0))

	orig := NewDeck()
	first := orig.cards[40]
	second := orig.cards[35]

	shuffled := NewDeck(Shuffle)
	if shuffled.cards[0] != first {
		t.Errorf("Expected the first card to be %s, received %s.", first, shuffled.cards[0])
	}
	if shuffled.cards[1] != second {
		t.Errorf("Expected the first card to be %s, received %s.", second, shuffled.cards[1])
	}
}

func TestJokers(t *testing.T) {
	deck := NewDeck(Jokers(4))
	count := 0
	for _, c := range deck.cards {
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
	for _, c := range deck.cards {
		if c.Rank == Two || c.Rank == Three {
			t.Error("Expected all twos and threes to be filtered out.")
		}
	}
}

func TestPacks(t *testing.T) {
	deck := NewDeck(Packs(3))
	// 13 ranks * 4 suits * 3 decks
	if len(deck.cards) != 13*4*3 {
		t.Errorf("Expected %d cards, received %d cards.", 13*4*3, len(deck.cards))
	}
}
