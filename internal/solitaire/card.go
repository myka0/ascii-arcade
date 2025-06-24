package solitaire

import (
	"strings"
)

var (
	values = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	suits  = []string{"♠", "♣", "♦", "♥"}
)

const (
	Spade = iota
	Diamond
	Heart
	Club
)

const (
	Ace = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

// Card represents a single playing card with a suit, rank, and visual states.
type Card struct {
	Suit       int
	Rank       int
	FaceDown   bool
	IsSelected bool
}

// NewCard creates a new card with a specified suit and rank.
func NewCard(suit, rank int) Card {
	return Card{
		Suit:       suit,
		Rank:       rank,
		FaceDown:   false,
		IsSelected: false,
	}
}

// Flip toggles the FaceDown status of the card
func (c *Card) Flip() {
	c.FaceDown = !c.FaceDown
}

// View returns a string representation of the card.
func (c Card) View() string {
	var card string
	if c.FaceDown {
		card = ViewFaceDown()
	} else {
		card = c.ViewCard()
	}

	return card
}

// ViewFaceDown renders a stylized representation of the back of a card.
func ViewFaceDown() string {
	style := FGWhite.Render
	return strings.Join([]string{
		style("╭─────╮"),
		style("│╱╱╱╱╱│"),
		style("│╱╱╱╱╱│"),
		style("│╱╱╱╱╱│"),
		style("│╱╱╱╱╱│"),
		style("╰─────╯"),
	}, "\n")
}

// ViewEmptyCard renders an empty card shaped placeholder.
func ViewEmptyCard() string {
	style := FGEmpty.Render
	return strings.Join([]string{
		style("╭─────╮"),
		style("│     │"),
		style("│     │"),
		style("│     │"),
		style("│     │"),
		style("╰─────╯"),
	}, "\n")
}

// ViewTop renders only the top line of a card for expanded decks.
func (c Card) ViewTop() string {
	rank := values[c.Rank]
	suit := c.Suit

	// Choose style based on suit
	suitStyle := FGWhite
	if suit == 2 || suit == 3 { // ♥♦
		suitStyle = FGRed
	}

	// Render styled rank and suit
	r := suitStyle.Render(rank)
	s := suitStyle.Render(suits[suit])

	// Choose style for the border of the card
	style := FGWhite.Render
	if c.IsSelected {
		style = FGSelected.Render
	}

	// Render face down top
	if c.FaceDown {
		return style("╭─────╮")
	}

	// Render top with rank and suit
	line := style(strings.Repeat("─", 4-len(rank)))
	return style("╭") + r + s + line + style("╮")
}

// ViewCard renders the full face up card, including its rank, suit, and interior.
func (c Card) ViewCard() string {
	rank := values[c.Rank]
	suit := c.Suit

	// Choose style based on suit
	suitStyle := FGWhite
	if suit == 2 || suit == 3 { // ♥ ♦
		suitStyle = FGRed
	}

	// Render styled rank and suit
	r := suitStyle.Render(rank)
	s := suitStyle.Render(suits[suit])

	// Style for the face of the card
	face := suitStyle.Render

	// Choose style for the border of the card
	style := FGWhite.Render
	if c.IsSelected {
		style = FGSelected.Render
	}

	// Top and bottom borders with rank and suit
	line := style(strings.Repeat("─", 4-len(rank)))
	top := style("╭") + r + s + line + style("╮")
	bot := style("╰") + line + r + s + style("╯")

	lines := []string{top}

	switch rank {
	case "2":
		lines = append(lines,
			style("│     │"),
			style("│  ")+s+style("  │"),
			style("│  ")+s+style("  │"),
			style("│     │"),
		)
	case "3":
		lines = append(lines,
			style("│     │"),
			style("│  ")+s+style("  │"),
			style("│ ")+s+style(" ")+s+style(" │"),
			style("│     │"),
		)
	case "4":
		lines = append(lines,
			style("│     │"),
			style("│ ")+s+style(" ")+s+style(" │"),
			style("│ ")+s+style(" ")+s+style(" │"),
			style("│     │"),
		)
	case "5":
		lines = append(lines,
			style("│  ")+s+style("  │"),
			style("│")+s+style("   ")+s+style("│"),
			style("│  ")+s+style("  │"),
			style("│  ")+s+style("  │"),
		)
	case "6":
		lines = append(lines,
			style("│")+s+style("   ")+s+style("│"),
			style("│  ")+s+style("  │"),
			style("│  ")+s+style("  │"),
			style("│")+s+style("   ")+s+style("│"),
		)
	case "7":
		lines = append(lines,
			style("│  ")+s+style("  │"),
			style("│ ")+s+style(" ")+s+style(" │"),
			style("│ ")+s+style(" ")+s+style(" │"),
			style("│ ")+s+style(" ")+s+style(" │"),
		)
	case "8":
		lines = append(lines,
			style("│ ")+s+style(" ")+s+style(" │"),
			style("│ ")+s+style(" ")+s+style(" │"),
			style("│ ")+s+style(" ")+s+style(" │"),
			style("│ ")+s+style(" ")+s+style(" │"),
		)
	case "9":
		lines = append(lines,
			style("│ ")+s+style(" ")+s+style(" │"),
			style("│")+s+style(" ")+s+style(" ")+s+style("│"),
			style("│ ")+s+style(" ")+s+style(" │"),
			style("│ ")+s+style(" ")+s+style(" │"),
		)
	case "10":
		lines = append(lines,
			style("│")+s+style(" ")+s+style(" ")+s+style("│"),
			style("│ ")+s+style(" ")+s+style(" │"),
			style("│ ")+s+style(" ")+s+style(" │"),
			style("│")+s+style(" ")+s+style(" ")+s+style("│"),
		)
	case "J":
		lines = append(lines,
			style("│   ")+face("▗▖")+style("│"),
			style("│   ")+face("▐▌")+style("│"),
			style("│   ")+face("▐▌")+style("│"),
			style("│")+face("▝▚▄▞▘")+style("│"),
		)
	case "Q":
		lines = append(lines,
			style("│ ")+face("▄▄▄")+style(" │"),
			style("│")+face("▐▌ ▐▌")+style("│"),
			style("│")+face("▐▌ ▐▌")+style("│"),
			style("│ ")+face("▀▀▜▖")+style("│"),
		)
	case "K":
		lines = append(lines,
			style("│")+face("▗▖ ▗▖")+style("│"),
			style("│")+face("▐▌▗▞▘")+style("│"),
			style("│")+face("▐▛▚▖ ")+style("│"),
			style("│")+face("▐▌ ▐▌")+style("│"),
		)
	case "A":
		lines = append(lines,
			style("│ ")+face("▗▄▖")+style(" │"),
			style("│")+face("▐▌ ▐▌")+style("│"),
			style("│")+face("▐▛▀▜▌")+style("│"),
			style("│")+face("▐▌ ▐▌")+style("│"),
		)
	}

	// Add bottom border
	lines = append(lines, bot)
	return strings.Join(lines, "\n")
}
