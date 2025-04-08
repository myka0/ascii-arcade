package wordle

import (
  tea "github.com/charmbracelet/bubbletea"
)

type wordleModel struct {
}

func InitWordleModel() wordleModel {
  m := wordleModel {
  }

  return m
}

func (m wordleModel) Init() tea.Cmd {
  return nil
}

func (m wordleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
    case "ctrl+c":
      return m, tea.Quit
    }
  }

  return m, nil
}

func (m wordleModel) View() string {
  var output string

  output += "Not implemented yet. Press ctrl+c to quit."
  return output
}

