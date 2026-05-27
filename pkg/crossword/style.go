package crossword

import (
	"ascii-arcade/internal/colors"

	"charm.land/lipgloss/v2"
)

const (
	Padding       = 2
	ClueWidth     = 40
	FullClueWidth = ClueWidth + Padding

	LowerBar = "▄▄▄▄▄"
	UpperBar = "▀▀▀▀▀"
	FullBar  = "█████"
	Blank    = "     "
)

var (
	Across    = colors.Purple
	Down      = colors.Pink
	Cursor    = colors.Yellow
	Incorrect = colors.Orange
	Even      = colors.Light1
	Odd       = colors.Light2
	DarkText  = colors.Dark1
	LightText = colors.Light2
	GreyText  = colors.Medium2
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

	// Clue Styles
	AcrossClue = lipgloss.NewStyle().
			Width(FullClueWidth).
			Foreground(DarkText).
			Background(Across).
			Padding(0, 1).
			Bold(true)

	DownClue = lipgloss.NewStyle().
			Width(FullClueWidth).
			Foreground(DarkText).
			Background(Down).
			Padding(0, 1).
			Bold(true)

	SolvedClue = lipgloss.NewStyle().
			Width(FullClueWidth).
			Foreground(GreyText).
			Padding(0, 1).
			Bold(true)

	NormalClue = lipgloss.NewStyle().
			Width(FullClueWidth).
			Foreground(LightText).
			Padding(0, 1).
			Bold(true)

	// Solid bars used in viewMargin / viewTopRow
	CursorLowerBar = FGCursor.Render(LowerBar)
	CursorUpperBar = FGCursor.Render(UpperBar)

	IncorrectLowerBar = FGIncorrect.Render(LowerBar)
	IncorrectUpperBar = FGIncorrect.Render(UpperBar)
	IncorrectFullBar  = FGIncorrect.Render(FullBar)

	AcrossLowerBar = FGAcross.Render(LowerBar)
	AcrossUpperBar = FGAcross.Render(UpperBar)

	DownLowerBar = FGDown.Render(LowerBar)
	DownUpperBar = FGDown.Render(UpperBar)
	DownFullBar  = FGDown.Render(FullBar)

	// Bars rendered with combined fg/bg styles (top connector cells)
	IncorrectTopCursorLowerBar = IncorrectTopCursor.Render(LowerBar)
	IncorrectTopCursorUpperBar = IncorrectTopCursor.Render(UpperBar)

	CursorDownTopLowerBar = CursorDownTop.Render(LowerBar)
	CursorDownTopUpperBar = CursorDownTop.Render(UpperBar)

	IncorrectTopAcrossLowerBar = IncorrectTopAcross.Render(LowerBar)
	IncorrectTopAcrossUpperBar = IncorrectTopAcross.Render(UpperBar)

	IncorrectTopDownLowerBar = IncorrectTopDown.Render(LowerBar)
	IncorrectTopDownUpperBar = IncorrectTopDown.Render(UpperBar)

	// Parity keyed bars
	CursorTopLowerBar = [2]string{
		CursorTopEven.Render(LowerBar),
		CursorTopOdd.Render(LowerBar),
	}
	CursorTopUpperBar = [2]string{
		CursorTopEven.Render(UpperBar),
		CursorTopOdd.Render(UpperBar),
	}

	IncorrectTopLowerBar = [2]string{
		IncorrectTopEven.Render(LowerBar),
		IncorrectTopOdd.Render(LowerBar),
	}
	IncorrectTopUpperBar = [2]string{
		IncorrectTopEven.Render(UpperBar),
		IncorrectTopOdd.Render(UpperBar),
	}

	AcrossTopLowerBar = [2]string{
		AcrossTopEven.Render(LowerBar),
		AcrossTopOdd.Render(LowerBar),
	}
	AcrossTopUpperBar = [2]string{
		AcrossTopEven.Render(UpperBar),
		AcrossTopOdd.Render(UpperBar),
	}

	DownTopLowerBar = [2]string{
		DownTopEven.Render(LowerBar),
		DownTopOdd.Render(LowerBar),
	}
	DownTopUpperBar = [2]string{
		DownTopEven.Render(UpperBar),
		DownTopOdd.Render(UpperBar),
	}

	EmptyBGLowerBar = [2]string{
		FGEven.Render(LowerBar),
		FGOdd.Render(LowerBar),
	}
	EmptyBGUpperBar = [2]string{
		FGEven.Render(UpperBar),
		FGOdd.Render(UpperBar),
	}

	TopLowerBar = [2]string{
		TopEven.Render(LowerBar),
		TopOdd.Render(LowerBar),
	}
)
