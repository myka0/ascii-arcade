package wordle

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const filename = "data/wordle/games.db"

// GetLatestDate returns the date of the most recent Wordle game state.
func GetLatestDate() (string, error) {
	// Open the database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return "", fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Get and return the latest date
	var date string
	err = db.QueryRow(`SELECT date FROM wordle ORDER BY date DESC LIMIT 1`).Scan(&date)
	return date, err
}

// SaveToFile writes the current game state to a SQLite database.
func (m *WordleModel) SaveToFile() error {
	// Open the database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS wordle (
			date TEXT PRIMARY KEY,
			answer TEXT,
			guesses TEXT,
			cursor_x INTEGER,
			cursor_y INTEGER,
			keyboard TEXT
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	// Convert byte slices to strings for JSON serialization
	var guesses [6]string
	for i, g := range m.guesses {
		guesses[i] = string(g[:])
	}

	// Convert byte keys to string keys for JSON serialization
	keyboard := make(map[string]int)
	for k, v := range m.keyboard {
		keyboard[string(k)] = v
	}

	// Convert slices to JSON
	guessesJSON, _ := json.Marshal(guesses)
	keyboardJSON, _ := json.Marshal(keyboard)

	// Insert the data into the database
	_, err = db.Exec(`
		INSERT OR REPLACE INTO wordle (date, answer, guesses, cursor_x, cursor_y, keyboard)
		VALUES (?, ?, ?, ?, ?, ?)
	`, m.date, string(m.answer[:]), guessesJSON, m.cursorX, m.cursorY, keyboardJSON)

	return err
}

// LoadFromFile loads a wordle game state from the SQLite database.
func LoadFromFile(date string) (WordleModel, error) {
	model := WordleModel{
		date:     date,
		answer:   [5]byte{},
		guesses:  [6][5]byte{},
		cursorX:  0,
		cursorY:  0,
		keyboard: make(map[byte]int, 26),
		message:  "",
	}

	// Open the database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return WordleModel{}, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Get the saved game data from the database
	row := db.QueryRow(`SELECT answer, guesses, cursor_x, cursor_y, keyboard FROM wordle WHERE date = ?`, date)
	var answer, guessesJSON, keyboardJSON []byte
	if err := row.Scan(&answer, &guessesJSON, &model.cursorX, &model.cursorY, &keyboardJSON); err != nil {
		return WordleModel{}, err
	}

	// Load answer into fixed-size byte array
	copy(model.answer[:], answer)

	// Decode and copy guesses into fixed-size byte arrays
	var guesses [6]string
	json.Unmarshal([]byte(guessesJSON), &guesses)
	for i, guess := range guesses {
		copy(model.guesses[i][:], guess)
	}

	// Decode and map keyboard state
	var keyboard map[string]int
	json.Unmarshal(keyboardJSON, &keyboard)
	for k, v := range keyboard {
		model.keyboard[k[0]] = v
	}

	return model, nil
}
