package connections

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const filename = "data/connections/games.db"

// GetLatestDate returns the date of the most recent connections game state.
func GetLatestDate() (string, error) {
	// Open the database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return "", fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Get and return the latest date
	var date string
	err = db.QueryRow(`SELECT date FROM connections ORDER BY date DESC LIMIT 1`).Scan(&date)
	return date, err
}

// SaveToFile persists the current connections game state to a SQLite database.
func (m *ConnectionsModel) SaveToFile() error {
	// Open the database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS connections (
			date TEXT PRIMARY KEY,
			word_groups TEXT,
			guess_history TEXT,
			revealed_word_groups TEXT,
			mistakes_remaining INTEGER
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	// Convert slices to JSON
	wordGroupsJSON, _ := json.Marshal(m.wordGroups)
	guessHistoryJSON, _ := json.Marshal(m.guessHistory)
	revealedWordGroupsJSON, _ := json.Marshal(m.revealedWordGroups)

	// Insert the data into the database
	_, err = db.Exec(`
		INSERT OR REPLACE INTO connections (date, word_groups, guess_history, revealed_word_groups, mistakes_remaining)
		VALUES (?, ?, ?, ?, ?)
	`, m.date, wordGroupsJSON, guessHistoryJSON, revealedWordGroupsJSON, m.mistakesRemaining)

	return err
}

// LoadFromFile loads a connections game state from the SQLite database.
func LoadFromFile(date string) (ConnectionsModel, error) {
	model := ConnectionsModel{}

	// Open the database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return model, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Get the saved game data from the database
	row := db.QueryRow(`SELECT word_groups, guess_history, revealed_word_groups, mistakes_remaining FROM connections WHERE date = ?`, date)
	var wordGroupsJSON, guessHistoryJSON, revealedWordGroupsJSON []byte
	var mistakesRemaining int
	if err := row.Scan(&wordGroupsJSON, &guessHistoryJSON, &revealedWordGroupsJSON, &mistakesRemaining); err != nil {
		return model, err
	}

	// Convert JSON to slices
	var wordGroups [4]WordGroup
	var guessHistory [][]string
	var revealedWordGroups [][]string
	json.Unmarshal(wordGroupsJSON, &wordGroups)
	json.Unmarshal(guessHistoryJSON, &guessHistory)
	json.Unmarshal(revealedWordGroupsJSON, &revealedWordGroups)

	// Set up the model with data from the saved game
	model.date = date
	model.wordGroups = wordGroups
	model.guessHistory = guessHistory
	model.revealedWordGroups = revealedWordGroups
	model.mistakesRemaining = mistakesRemaining

	return model, nil
}
