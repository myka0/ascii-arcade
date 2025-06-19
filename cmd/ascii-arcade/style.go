package main

import (
	"ascii-arcade/internal/colors"

	"github.com/charmbracelet/lipgloss/v2"
)

var (
	MenuWidth        = 48
	KeyBindMenuWidth = 16

	Text = lipgloss.NewStyle().Foreground(colors.Light2)

	CenteredText = Text.
			Width(MenuWidth).
			Align(lipgloss.Center)

	// List styles
	ListStyle = lipgloss.NewStyle().
			MarginLeft(16)

	ListHeader = lipgloss.NewStyle().
			Foreground(colors.Dark).
			Background(colors.Purple).
			Padding(0, 1).
			MarginTop(1).
			Bold(true)

	ListEntry = lipgloss.NewStyle().
			Foreground(colors.Light2).
			MarginLeft(2)

	SelectedListEntry = lipgloss.NewStyle().
				Foreground(colors.Pink)

	// Menu styles
	MenuStyle = lipgloss.NewStyle().
			Width(MenuWidth)

	KeyBindMenu = lipgloss.NewStyle().
			Width(KeyBindMenuWidth)

	// Keybinding styles
	KeyStyle = lipgloss.NewStyle().
			Foreground(colors.Medium1)

	KeyActionStyle = lipgloss.NewStyle().
			Foreground(colors.Medium2)

	Header = lipgloss.NewStyle().Foreground(colors.Purple).Render(`
         █████╗ ███████╗ ██████╗██╗██╗
        ██╔══██╗██╔════╝██╔════╝██║██║
        ███████║███████╗██║     ██║██║
        ██╔══██║╚════██║██║     ██║██║
        ██║  ██║███████║╚██████╗██║██║
        ╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝╚═╝

 █████╗ ██████╗  ██████╗ █████╗ ██████╗ ███████╗
██╔══██╗██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔════╝
███████║██████╔╝██║     ███████║██║  ██║█████╗
██╔══██║██╔══██╗██║     ██╔══██║██║  ██║██╔══╝
██║  ██║██║  ██║╚██████╗██║  ██║██████╔╝███████╗
╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚═════╝ ╚══════╝`,
	)
)
