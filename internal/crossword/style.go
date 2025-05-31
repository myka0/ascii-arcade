package crossword

import (
	"crossword/internal/colors"
	"github.com/charmbracelet/lipgloss"
)

var (
	Across    = colors.Purple
	Down      = colors.Pink
	Cursor    = colors.Yellow
	Incorrect = colors.Orange
	Even      = colors.Light1
	Odd       = colors.Light2
	DarkText  = colors.Dark
	LightText = colors.Light1
	GreyText  = colors.Medium
	Border    = colors.Blue

	FGAcross    = lipgloss.NewStyle().Foreground(Across)
	FGDown      = lipgloss.NewStyle().Foreground(Down)
	FGCursor    = lipgloss.NewStyle().Foreground(Cursor)
	FGIncorrect = lipgloss.NewStyle().Foreground(Incorrect)
	FGEven      = lipgloss.NewStyle().Foreground(Even)
	FGOdd       = lipgloss.NewStyle().Foreground(Odd)
	FGLightText = lipgloss.NewStyle().Foreground(LightText)
	FGGreyText  = lipgloss.NewStyle().Foreground(GreyText)
	FGBorder    = lipgloss.NewStyle().Foreground(Border)

	BorderStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Border)

	// General Cells
	CursorCell = lipgloss.NewStyle().
			Foreground(DarkText).
			Background(Cursor).
			Bold(true)

	AcrossCell = lipgloss.NewStyle().
			Foreground(DarkText).
			Background(Across).
			Bold(true)

	DownCell = lipgloss.NewStyle().
			Foreground(DarkText).
			Background(Down).
			Bold(true)

	IncorrectCell = lipgloss.NewStyle().
			Foreground(DarkText).
			Background(Incorrect).
			Bold(true)

	EvenCell = lipgloss.NewStyle().
			Foreground(DarkText).
			Background(Even).
			Bold(true)

	OddCell = lipgloss.NewStyle().
		Foreground(DarkText).
		Background(Odd).
		Bold(true)

	// Cursor Top
	CursorTopEven = lipgloss.NewStyle().
			Foreground(Cursor).
			Background(Even)

	CursorTopOdd = lipgloss.NewStyle().
			Foreground(Cursor).
			Background(Odd)

	CursorDownTop = lipgloss.NewStyle().
			Foreground(Cursor).
			Background(Down)

	// Across Top
	AcrossTopEven = lipgloss.NewStyle().
			Foreground(Across).
			Background(Even)

	AcrossTopOdd = lipgloss.NewStyle().
			Foreground(Across).
			Background(Odd)

	// Down Top
	DownTopEven = lipgloss.NewStyle().
			Foreground(Down).
			Background(Even)

	DownTopOdd = lipgloss.NewStyle().
			Foreground(Down).
			Background(Odd)

	// Incorrect Top
	IncorrectTopCursor = lipgloss.NewStyle().
				Foreground(Incorrect).
				Background(Cursor)

	IncorrectTopAcross = lipgloss.NewStyle().
				Foreground(Incorrect).
				Background(Across)

	IncorrectTopDown = lipgloss.NewStyle().
				Foreground(Incorrect).
				Background(Down)

	IncorrectTopEven = lipgloss.NewStyle().
				Foreground(Incorrect).
				Background(Even)

	IncorrectTopOdd = lipgloss.NewStyle().
			Foreground(Incorrect).
			Background(Odd)

	// General Top
	TopEven = lipgloss.NewStyle().
		Foreground(Even).
		Background(Odd)

	TopOdd = lipgloss.NewStyle().
		Foreground(Odd).
		Background(Even)
)
