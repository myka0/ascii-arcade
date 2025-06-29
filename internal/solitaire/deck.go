package solitaire

import (
	"hash/maphash"
	"math/rand"
)

type Deck struct {
	Cards      []*Card
	isExpanded bool
}

// NewDeck creates a new deck of cards from a given slice of cards.
func NewDeck(cards []*Card) Deck {
	return Deck{
		Cards: cards,
	}
}

// NewFullDeck creates a completely new deck of cards.
func NewFullDeck() Deck {
	cards := make([]*Card, 52)
	for i := range 52 {
		cards[i] = NewCard(i/13, i%13)
	}
	return NewDeck(cards)
}

// NewEmptyDeck creates an empty deck of cards.
func NewEmptyDeck() Deck {
	return NewDeck(make([]*Card, 0))
}

// View returns a string representation of the deck.
func (d Deck) View() string {
	// If the deck is empty, return an empty card
	if d.Size() == 0 {
		return ViewEmptyCard()
	}

	// Expand cards
	if d.isExpanded {
		var view string
		for i := range d.Size() - 1 {
			view += d.Cards[i].ViewTop() + "\n"
		}
		return view + d.Top().View()
	}

	// View the top card
	return d.Top().View()
}

// Add add new cards to the deck.
func (d *Deck) Add(cards ...*Card) {
	d.Cards = append(d.Cards, cards...)
}

// Get returns the card at the given index.
func (d Deck) Get(idx int) *Card {
	return d.Cards[idx]
}

// Top returns the top card of the deck.
func (d Deck) Top() *Card {
	return d.Get(d.Size() - 1)
}

// Pop removes and returns the top card of the deck.
func (d *Deck) Pop() *Card {
	card := d.Top()
	d.Cards = d.Cards[:d.Size()-1]
	return card
}

// Size returns the number of cards in the deck.
func (d Deck) Size() int {
	return len(d.Cards)
}

// Expand expands the deck.
func (d *Deck) Expand() {
	d.isExpanded = true
}

// Shuffle shuffles the deck.
func (d Deck) Shuffle() {
	// Use the maphash package to generate a random seed
	generator := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
	generator.Shuffle(d.Size(), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}
