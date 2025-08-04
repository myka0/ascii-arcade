package wordle

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const filename = "data/wordle/games.db"

// JSON response object from NYT
type WordleResponse struct {
	ID              int    `json:"id"`
	Solution        string `json:"solution"`
	PrintDate       string `json:"print_date"`
	DaysSinceLaunch int    `json:"days_since_launch"`
	Editor          string `json:"editor"`
}

// LoadGame returns the Wordle game state for a given date.
func LoadGame(date string) (WordleModel, error) {
	// Try loading the saved game from the database
	model, err := LoadFromFile(date)
	if err == nil && model.answer != [5]byte{} {
		return model, nil
	}

	// Initialize a new empty game state
	model = WordleModel{
		date:     date,
		guesses:  [6][5]byte{},
		keyboard: make(map[byte]int, 26),
	}
	model.handleReset()

	// Try fetching from the web if not in database
	answer, fetchErr := fetchWordleAnswer(date)
	if fetchErr != nil {
		return model, fetchErr
	}

	model.answer = answer
	return model, nil
}

// fetchWordleAnswer fetches the Wordle answer for today from the NYT API.
func fetchWordleAnswer(date string) ([5]byte, error) {
	var answer [5]byte
	url := fmt.Sprintf("https://www.nytimes.com/svc/wordle/v2/%s.json", date)

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		return answer, fmt.Errorf("error fetching Wordle answer: %v", err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return answer, fmt.Errorf("non-OK HTTP status: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return answer, fmt.Errorf("error reading response body: %v", err)
	}

	// Decode JSON
	var wordle WordleResponse
	if err := json.Unmarshal(body, &wordle); err != nil {
		return answer, fmt.Errorf("error decoding JSON: %v\nbody: %s", err, string(body))
	}

	copy(answer[:], strings.ToUpper(wordle.Solution))
	return answer, nil
}

// SaveToFile writes the current game state to a SQLite database.
func (m *WordleModel) SaveToFile() error {
	// Create the data directory if it doesn't exist
	os.MkdirAll("data/wordle", 0755)

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
		keyboard: make(map[byte]int, 26),
	}

	model.handleReset()

	// Open the database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return model, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Get the saved game data from the database
	row := db.QueryRow(`SELECT answer, guesses, cursor_x, cursor_y, keyboard FROM wordle WHERE date = ?`, date)
	var answer, guessesJSON, keyboardJSON []byte
	if err := row.Scan(&answer, &guessesJSON, &model.cursorX, &model.cursorY, &keyboardJSON); err != nil {
		return model, err
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
