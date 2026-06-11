package minesweeper

import (
	"ascii-arcade/internal/colors"

	"charm.land/lipgloss/v2"
)

const (
	Width = 5

	LowerBar = "▄▄▄▄▄"
	UpperBar = "▀▀▀▀▀"

	Mine    = "\U000F0691"
	Flag    = "⚑"
	MineHit = "\U000F0691"
)

var (
	Even     = colors.Light1
	Odd      = colors.Light2
	Revealed = colors.Medium2
	Cursor   = colors.Blue
	Hit      = colors.Red

	FGEven     = lipgloss.NewStyle().Foreground(Even)
	FGOdd      = lipgloss.NewStyle().Foreground(Odd)
	FGRevealed = lipgloss.NewStyle().Foreground(Revealed)
	FGCursor   = lipgloss.NewStyle().Foreground(Cursor)
	FGHit      = lipgloss.NewStyle().Foreground(Hit)

	LabelStyle = lipgloss.NewStyle().
			Foreground(colors.Dark1).
			Background(colors.Purple).
			Padding(0, 1).
			MarginBottom(1).
			Bold(true)

	ButtonStyle = lipgloss.NewStyle().
			Foreground(colors.Dark1).
			Background(colors.Purple).
			Padding(0, 1).
			MarginTop(1).
			Bold(true)

	ListEntry = lipgloss.NewStyle().
			Foreground(colors.Light2).
			MarginLeft(2)

	SelectedListEntry = lipgloss.NewStyle().
				Foreground(colors.Pink)

	MessageStyle = lipgloss.NewStyle().
			Foreground(colors.Light2).
			Bold(true)

	OddCell = lipgloss.NewStyle().
		Align(lipgloss.Center).
		Foreground(colors.Dark1).
		Background(Odd).
		Width(Width).
		Bold(true)

	EvenCell = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(colors.Dark1).
			Background(Even).
			Width(Width).
			Bold(true)

	RevealedCell = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Background(Revealed).
			Width(Width).
			Bold(true)

	CursorCell = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(colors.Dark1).
			Background(Cursor).
			Width(Width).
			Bold(true)

	HitCell = lipgloss.NewStyle().
		Align(lipgloss.Center).
		Foreground(colors.Dark1).
		Background(Hit).
		Width(Width).
		Render(MineHit)

	MineCell = RevealedCell.Foreground(colors.Dark1).Render(Mine)

	CursorMarginEven = lipgloss.NewStyle().
				Foreground(Cursor).
				Background(Even)

	CursorMarginOdd = lipgloss.NewStyle().
			Foreground(Cursor).
			Background(Odd)

	CursorRevealedMargin = lipgloss.NewStyle().
				Foreground(Cursor).
				Background(Revealed)

	HitMarginEven = lipgloss.NewStyle().
			Foreground(Hit).
			Background(Even)

	HitMarginOdd = lipgloss.NewStyle().
			Foreground(Hit).
			Background(Odd)

	HitRevealedMargin = lipgloss.NewStyle().
				Foreground(Hit).
				Background(Revealed)

	RevealedMarginEven = lipgloss.NewStyle().
				Foreground(Revealed).
				Background(Even)

	RevealedMarginOdd = lipgloss.NewStyle().
				Foreground(Revealed).
				Background(Odd)

	MarginEven = lipgloss.NewStyle().
			Foreground(Even).
			Background(Odd)

	MarginOdd = lipgloss.NewStyle().
			Foreground(Odd).
			Background(Even)

	CursorLowerBar   = FGCursor.Render(LowerBar)
	CursorUpperBar   = FGCursor.Render(UpperBar)
	CursorFullBar    = CursorCell.Render("")
	CursorFlaggedBar = CursorCell.Render(Flag)

	HitLowerBar = FGHit.Render(LowerBar)
	HitUpperBar = FGHit.Render(UpperBar)

	RevealedLowerBar = FGRevealed.Render(LowerBar)
	RevealedUpperBar = FGRevealed.Render(UpperBar)
	RevealedFullBar  = RevealedCell.Render("")

	CursorRevealedMarginLowerBar = CursorRevealedMargin.Render(LowerBar)
	CursorRevealedMarginUpperBar = CursorRevealedMargin.Render(UpperBar)

	HitRevealedMarginLowerBar = HitRevealedMargin.Render(LowerBar)
	HitRevealedMarginUpperBar = HitRevealedMargin.Render(UpperBar)

	CursorMarginLowerBar = [2]string{
		CursorMarginEven.Render(LowerBar),
		CursorMarginOdd.Render(LowerBar),
	}
	CursorMarginUpperBar = [2]string{
		CursorMarginEven.Render(UpperBar),
		CursorMarginOdd.Render(UpperBar),
	}

	HitMarginLowerBar = [2]string{
		HitMarginEven.Render(LowerBar),
		HitMarginOdd.Render(LowerBar),
	}
	HitMarginUpperBar = [2]string{
		HitMarginEven.Render(UpperBar),
		HitMarginOdd.Render(UpperBar),
	}

	RevealedMarginLowerBar = [2]string{
		RevealedMarginEven.Render(LowerBar),
		RevealedMarginOdd.Render(LowerBar),
	}
	RevealedMarginUpperBar = [2]string{
		RevealedMarginEven.Render(UpperBar),
		RevealedMarginOdd.Render(UpperBar),
	}

	EmptyBGLowerBar = [2]string{
		FGEven.Render(LowerBar),
		FGOdd.Render(LowerBar),
	}
	EmptyBGUpperBar = [2]string{
		FGEven.Render(UpperBar),
		FGOdd.Render(UpperBar),
	}

	MarginLowerBar = [2]string{
		MarginEven.Render(LowerBar),
		MarginOdd.Render(LowerBar),
	}

	FlaggedCell = [2]string{
		EvenCell.Render(Flag),
		OddCell.Render(Flag),
	}

	WrongFlaggedCell = [2]string{
		EvenCell.Foreground(colors.Red).Render(Flag),
		OddCell.Foreground(colors.Red).Render(Flag),
	}

	NormalBar = [2]string{
		EvenCell.Render(""),
		OddCell.Render(""),
	}

	ColoredAdjacentStyles = []string{
		RevealedCell.Render(""),
		RevealedCell.Foreground(colors.Blue).Render("1"),
		RevealedCell.Foreground(colors.Green).Render("2"),
		RevealedCell.Foreground(colors.Red).Render("3"),
		RevealedCell.Foreground(colors.Purple).Render("4"),
		RevealedCell.Foreground(colors.Orange).Render("5"),
		RevealedCell.Foreground(colors.Cyan).Render("6"),
		RevealedCell.Foreground(colors.Yellow).Render("7"),
		RevealedCell.Foreground(colors.Dark2).Render("8"),
	}

	CursorAdjacentStyles = []string{
		CursorCell.Render(""),
		CursorCell.Render("1"),
		CursorCell.Render("2"),
		CursorCell.Render("3"),
		CursorCell.Render("4"),
		CursorCell.Render("5"),
		CursorCell.Render("6"),
		CursorCell.Render("7"),
		CursorCell.Render("8"),
	}
)