package connections

import (
  tea "github.com/charmbracelet/bubbletea"
)

type connectionsModel struct {
}

func InitConnectionsModel() connectionsModel {
  m := connectionsModel {
  }

  return m
}

func (m connectionsModel) Init() tea.Cmd {
  return nil
}

func (m connectionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
    case "ctrl+c":
      return m, tea.Quit
    }
  }

  return m, nil
}

func (m connectionsModel) View() string {
  var output string

  output += "Not implemented yet. Press ctrl+c to quit."
  return output
}

