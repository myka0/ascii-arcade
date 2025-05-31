package crossword

import (
	"encoding/json"
	"fmt"
	"os"
)

const filename = "data/crossword/games.json"

// saveData represents the structure used to serialize crossword puzzle data to JSON.
type saveData struct {
	Across []string   `json:"across"`
	Down   []string   `json:"down"`
	Answer [15]string `json:"answer"`
	Grid   [15]string `json:"grid"`
}

// SaveToFile persists the current crossword puzzle state to a JSON file.
func (m *CrosswordModel) SaveToFile() error {
	// Read existing data from the file
	savedGames := make(map[string]saveData)
	fileData, err := os.ReadFile(filename)
	if err == nil {
		if err := json.Unmarshal(fileData, &savedGames); err != nil {
			return fmt.Errorf("error parsing existing data: %v", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Convert byte arrays to strings for JSON serialization
	var answer [15]string
	var grid [15]string
	for i := range 15 {
		answer[i] = string(m.answer[i][:])
		grid[i] = string(m.grid[i][:])
	}

	// Create or update entry for current date
	savedGames[m.date] = saveData{
		Across: m.acrossClues,
		Down:   m.downClues,
		Answer: answer,
		Grid:   grid,
	}

	// Write updated data back to file
	jsonData, err := json.MarshalIndent(savedGames, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}

	return os.WriteFile(filename, jsonData, 0644)
}

// LoadFromFile loads the most recent crossword puzzle state from the saved JSON file.
func LoadFromFile() (CrosswordModel, error) {
	model := CrosswordModel{}

	// Initialize the grid with empty spaces
	for i := range model.grid {
		model.grid[i] = [15]byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	}

	// Read the saved game data from file
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return model, fmt.Errorf("error reading file: %w", err)
	}

	// Parse the JSON data
	var data map[string]saveData
	if err := json.Unmarshal(fileContent, &data); err != nil {
		return model, fmt.Errorf("error parsing JSON: %w", err)
	}

	// Extract dates from the saved games and get the most recent entry
	dates := make([]string, 0, len(data))
	for date := range data {
		dates = append(dates, date)
	}
	entry := data[dates[0]]

	// Set up the model with data from the saved game
	model.date = dates[0]
	model.acrossClues = entry.Across
	model.downClues = entry.Down
	model.isAcross = true
	model.message = ""

	// Convert string data back to byte arrays for the model
	for i, row := range entry.Answer {
		copy(model.answer[i][:], row)
	}

	for i, row := range entry.Grid {
		copy(model.grid[i][:], row)
	}

	// Initialize grid numbers, clue indices, and solved status
	model.initializeGrid()

	return model, nil
}

// initializeGrid sets up the grid numbers, clue indices, and tracks solved clues.
func (m *CrosswordModel) initializeGrid() {
	gridNumIdx := 0
	acrossClueIdx := -1
	downClueIdx := 0

	for i := range 225 {
		row := i / 15
		col := i % 15

		// If the answer is a black cell, mark the grid cell as black too
		if m.answer[row][col] == '.' {
			m.grid[row][col] = '.'
		}

		// Count correct cells and filled cells for progress tracking
		if m.grid[row][col] == m.answer[row][col] {
			m.correctCount++
		}
		if m.grid[row][col] != ' ' {
			m.filledCount++
		}

		// Skip further processing if the cell is a black square
		if m.grid[row][col] == '.' {
			m.clueIndices[row][col] = Position{-1, -1} // Mark as not part of any clue
			continue
		}

		// Check if this cell starts an across clue
		if col == 0 || m.grid[row][col-1] == '.' {
			gridNumIdx++
			m.gridNums[row][col] = gridNumIdx // Assign a grid number to this cell

			// Check if this across clue is already completely solved
			currentX := col
			isSolved := true
			for currentX < 15 && m.grid[row][currentX] != '.' {
				if m.grid[row][currentX] == ' ' {
					isSolved = false
					break
				}
				currentX++
			}

			m.isAcrossSolved = append(m.isAcrossSolved, isSolved)
			acrossClueIdx++
		}

		m.clueIndices[row][col].X = acrossClueIdx

		// Check if this cell starts a down clue
		if row == 0 || m.grid[row-1][col] == '.' {
			// If this cell doesn't already have a grid number assign one
			if m.gridNums[row][col] == 0 {
				gridNumIdx++
				m.gridNums[row][col] = gridNumIdx
			}

			// Check if this down clue is already completely solved
			currentY := row
			isSolved := true
			for currentY < 15 && m.grid[currentY][col] != '.' {
				// Set the down clue index for all cells in this clue
				m.clueIndices[currentY][col].Y = downClueIdx

				if m.grid[currentY][col] == ' ' {
					isSolved = false
				}
				currentY++
			}

			m.isDownSolved = append(m.isDownSolved, isSolved)
			downClueIdx++
		}
	}
}
