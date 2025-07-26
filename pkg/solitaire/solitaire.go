package solitaire

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// SolitaireModel represents the state of a solitaire game.
type SolitaireModel struct {
	stock       Deck
	waste       Deck
	foundations [4]Deck
	tableau     [7]Deck
}

// InitSolitaireModel creates and initializes a new solitaire model.
func InitSolitaireModel() *SolitaireModel {
	// Create a new shuffled deck for the stock
	stock := NewFullDeck()
	stock.Shuffle()

	// Create a new empty deck for the waste
	waste := NewEmptyDeck()

	// Initialize 4 empty foundations
	var foundations [4]Deck
	for i := range foundations {
		foundations[i] = NewEmptyDeck()
	}

	// Initialize 7 tableau columns
	var tableau [7]Deck
	for i := range tableau {
		tableau[i] = NewEmptyDeck()
	}

	// Deal cards to the tableau
	for i := range tableau {
		for range i + 1 {
			tableau[i].Add(stock.Pop())
		}

		tableau[i].Top().FlipFaceUp()
		tableau[i].Expand()
	}

	// Construct and return the game model
	m := SolitaireModel{
		stock:       stock,
		waste:       waste,
		foundations: foundations,
		tableau:     tableau,
	}

	return &m
}

// Init implements the Bubble Tea interface for initialization.
func (m SolitaireModel) Init() tea.Cmd {
	return nil
}

// Update handles keypress events and updates the model state accordingly.
func (m *SolitaireModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle keyboard input
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r":
			return InitSolitaireModel(), nil
		case "space":
			m.handleDrawFromStock()
		case "w":
			m.handleWasteAction()
		case "u":
			m.handleUndo()
		case "!": // Shift+1
			m.handleFoundationAction(Spade)
		case "@": // Shift+2
			m.handleFoundationAction(Club)
		case "#": // Shift+3
			m.handleFoundationAction(Heart)
		case "$": // Shift+4
			m.handleFoundationAction(Diamond)
		default:
			// Check if key is a number between 1 and 7
			if num, err := strconv.Atoi(msg.String()); err == nil && num >= 1 && num <= 7 {
				m.handleTableauAction(num - 1)
			}
		}

	// Handle mouse input
	case tea.MouseMsg:
		switch msg := msg.(type) {
		case tea.MouseClickMsg:
			m.handleMouseClick(msg)
		}
	}

	return m, nil
}

// handleDrawFromStock transfers a card from stock to waste.
func (m *SolitaireModel) handleDrawFromStock() {
	// If the stock is empty, recycle all waste cards into stock
	if m.stock.Size() == 0 {
		for range m.waste.Size() {
			card := m.waste.Pop()
			card.FlipFaceDown()
			m.stock.Add(card)
		}

		return
	}

	// Otherwise, draw a card from stock to waste and flip it face up
	card := m.stock.Pop()
	m.waste.Add(card)
	card.FlipFaceUp()
}

// handleTableauAction attempts to move cards from the specified tableau column.
func (m *SolitaireModel) handleTableauAction(col int) {
	source := &m.tableau[col]

	// If the selected tableau column is empty, do nothing
	if source.Size() == 0 {
		return
	}

	topCard := source.Top()

	// Try moving top card to the foundation
	if m.canMoveToFoundation(*topCard) {
		m.foundations[topCard.Suit].Add(source.Pop())

		// If there are cards left in the column, flip the new top card face up
		if source.Size() != 0 {
			source.Top().FlipFaceUp()
		}

		return
	}

	// Try moving a face up sequence to another tableau column
	for i, card := range source.Cards {
		if card.FaceDown {
			continue
		}

		// Try to move this sequence to a different column
		for targetCol := range m.tableau {
			// Skip if same column or if the card can't be moved
			if targetCol == col || !m.canMoveToTableau(*card, m.tableau[targetCol]) {
				continue
			}

			// Move the sequence of cards
			m.tableau[targetCol].Add(source.Cards[i:]...)

			// Remove them from the source column
			for range len(source.Cards) - i {
				source.Pop()
			}

			// Flip the new top card if any cards remain
			if source.Size() != 0 {
				source.Top().FlipFaceUp()
			}

			return // move only once
		}
	}
}

// handleWasteAction tries to move the top card from the waste pile.
func (m *SolitaireModel) handleWasteAction() {
	if m.waste.Size() == 0 {
		return // no cards to move
	}

	card := m.waste.Top()

	// Try to move to the foundation
	if m.canMoveToFoundation(*card) {
		m.foundations[card.Suit].Add(m.waste.Pop())
		return
	}

	// Try to move to a tableau column
	for i := range m.tableau {
		if m.canMoveToTableau(*card, m.tableau[i]) {
			m.tableau[i].Add(card)
			m.waste.Pop()
			break
		}
	}
}

// handleFoundationAction tries to move the top card from a specific foundation.
func (m *SolitaireModel) handleFoundationAction(suit int) {
	foundation := &m.foundations[suit]

	if foundation.Size() == 0 {
		return // no cards to move
	}

	card := foundation.Top()

	// Try to place the card on a valid tableau column
	for i := range m.tableau {
		if m.canMoveToTableau(*card, m.tableau[i]) {
			m.tableau[i].Add(foundation.Pop())
			break
		}
	}
}

// TODO: handleUndo
func (m *SolitaireModel) handleUndo() {

}

// TODO: handleMouseClick
func (m *SolitaireModel) handleMouseClick(msg tea.MouseMsg) {
}

// canMoveToFoundation checks if the given card can legally be placed onto its foundation pile.
func (m SolitaireModel) canMoveToFoundation(card Card) bool {
	foundation := m.foundations[card.Suit]

	// Only Aces can be placed on empty foundation piles
	if foundation.Size() == 0 {
		return card.Rank == Ace
	}

	// Card can be placed if it's the next rank up
	return card.Rank == foundation.Top().Rank+1
}

// canMoveToTableau checks if the given card can be placed on the target tableau pile.
func (m SolitaireModel) canMoveToTableau(card Card, tableau Deck) bool {
	// Only Kings can be placed on empty tableau piles
	if tableau.Size() == 0 {
		return card.Rank == King
	}

	// Cards must alternate colors
	cardIsBlack := card.Suit <= Club
	topIsBlack := tableau.Top().Suit <= Club
	if cardIsBlack == topIsBlack {
		return false
	}

	// Card can be placed if it's the next rank down
	return card.Rank == tableau.Top().Rank-1
}

// View renders the entire Solitaire board.
func (m SolitaireModel) View() string {
	// Render top row: Stock, Waste, and Foundations
	topRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.stock.View(),
		m.waste.View(),
		ViewCardSpacer(),
		m.foundations[0].View(),
		m.foundations[1].View(),
		m.foundations[2].View(),
		m.foundations[3].View(),
	)

	// Render middle row: Tableau column hints
	var columnHints []string
	for i := range 7 {
		columnHints = append(columnHints, TableauColumnHint.Render(strconv.Itoa(i+1)))
	}
	middleRow := lipgloss.JoinHorizontal(lipgloss.Top, columnHints...)

	// Render bottom row: Tableau columns
	bottomRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.tableau[0].View(),
		m.tableau[1].View(),
		m.tableau[2].View(),
		m.tableau[3].View(),
		m.tableau[4].View(),
		m.tableau[5].View(),
		m.tableau[6].View(),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		topRow,
		middleRow,
		bottomRow,
	)
}
