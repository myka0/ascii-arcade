package crossword

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const filename = "data/crossword/games.db"

// GetLatestDate returns the date of the most recent crossword puzzle state.
func GetLatestDate() (string, error) {
	// Open the database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return "", fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Get and return the latest date
	var date string
	err = db.QueryRow(`SELECT date FROM crosswords ORDER BY date DESC LIMIT 1`).Scan(&date)
	return date, err
}

// SaveToFile persists the current crossword puzzle state to a SQLite database.
func (m *CrosswordModel) SaveToFile() error {
	// Open the database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS crosswords (
			date TEXT PRIMARY KEY,
			across TEXT,
			down TEXT,
			answer TEXT,
			grid TEXT
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	// Convert byte slices to strings for JSON serialization
	var answer, grid [15]string
	for i := range 15 {
		answer[i] = string(m.answer[i][:])
		grid[i] = string(m.grid[i][:])
	}

	// Convert slices to JSON
	acrossJSON, _ := json.Marshal(m.acrossClues)
	downJSON, _ := json.Marshal(m.downClues)
	answerJSON, _ := json.Marshal(answer)
	gridJSON, _ := json.Marshal(grid)

	// Insert the data into the database
	_, err = db.Exec(`
		INSERT OR REPLACE INTO crosswords (date, across, down, answer, grid)
		VALUES (?, ?, ?, ?, ?)
	`, m.date, acrossJSON, downJSON, answerJSON, gridJSON)

	return err
}

// LoadFromFile loads a crossword puzzle state from the SQLite database.
func LoadFromFile(date string) (CrosswordModel, error) {
	model := CrosswordModel{}

	// Open the database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return model, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Get the saved game data from the database
	row := db.QueryRow(`SELECT across, down, answer, grid FROM crosswords WHERE date = ?`, date)
	var acrossJSON, downJSON, answerJSON, gridJSON []byte
	if err := row.Scan(&acrossJSON, &downJSON, &answerJSON, &gridJSON); err != nil {
		return model, err
	}

	// Convert JSON to slices
	var across, down []string
	var answer, grid [15]string
	json.Unmarshal(acrossJSON, &across)
	json.Unmarshal(downJSON, &down)
	json.Unmarshal(answerJSON, &answer)
	json.Unmarshal(gridJSON, &grid)

	// Set up the model with data from the saved game
	model.date = date
	model.acrossClues = across
	model.downClues = down
	model.isAcross = true
	model.message = ""

	// Convert string data back to byte arrays for the model
	for i := range 15 {
		copy(model.answer[i][:], answer[i])
		copy(model.grid[i][:], grid[i])
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
