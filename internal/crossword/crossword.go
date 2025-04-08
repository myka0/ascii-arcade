package crossword

import (
  tea "github.com/charmbracelet/bubbletea"
)

type crosswordModel struct {
}

func InitCrosswordModel() crosswordModel {
  m := crosswordModel {
  }

  return m
}

func (m crosswordModel) Init() tea.Cmd {
  return nil
}

func (m crosswordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
    case "ctrl+c":
      return m, tea.Quit
    }
  }

  return m, nil
}

func (m crosswordModel) View() string {
  var output string

  output += "Not implemented yet. Press ctrl+c to quit."
  return output
}
