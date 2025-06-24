package solitaire

type Deck struct {
	Cards      []Card
	isExpanded bool
}

// NewDeck creates a new deck of cards from a given slice of cards.
func NewDeck(cards []Card) Deck {
	return Deck{
		Cards: cards,
	}
}

// NewFullDeck creates a completely new deck of cards.
func NewFullDeck() Deck {
	cards := make([]Card, 52)
	for i := range 52 {
		cards[i] = NewCard(i/13, i%13)
	}
	return NewDeck(cards)
}

// View returns a string representation of the deck.
func (d Deck) View() string {
	// If the deck is empty, return an empty card
	if len(d.Cards) == 0 {
		return ViewEmptyCard()
	}

	// Expand cards
	if d.isExpanded {
		var view string
		for i := range len(d.Cards) - 1 {
			view += d.Cards[i].ViewTop() + "\n"
		}
		return view + d.Top().View()
	}

	// View the top card
	return d.Top().View()
}

// Add add new cards to the deck.
func (d *Deck) Add(cards ...Card) {
	d.Cards = append(d.Cards, cards...)
}

// Get returns the card at the given index.
func (d Deck) Get(idx int) Card {
	return d.Cards[idx]
}

// Top returns the top card of the deck.
func (d Deck) Top() Card {
	return d.Get(len(d.Cards) - 1)
}

// Pop removes and returns the top card of the deck.
func (d Deck) Pop() Card {
	return d.Get(0)
}
