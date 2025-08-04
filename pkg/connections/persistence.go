package connections

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const filename = "data/connections/games.db"

// JSON response objects from NYT
type ConnectionsResponse struct {
	Status     string     `json:"status"`
	ID         int        `json:"id"`
	PrintDate  string     `json:"print_date"`
	Editor     string     `json:"editor"`
	Categories []Category `json:"categories"`
}

type Category struct {
	Title string `json:"title"`
	Cards []Card `json:"cards"`
}

type Card struct {
	Content  string `json:"content"`
	Position int    `json:"position"`
}

// LoadGame returns the Wordle game state for a given date.
func LoadGame(date string) (ConnectionsModel, error) {
	// Try loading the saved game from the database
	model, err := LoadFromFile(date)
	if err == nil {
		return model, nil
	}

	// Initialize a new empty game state
	model = ConnectionsModel{
		date: date,
	}
	model.handleReset()

	// Try fetching from the web if not in database
	groups, fetchErr := fetchConnectionsGroups(date)
	if fetchErr != nil {
		return model, fetchErr
	}

	model.wordGroups = groups
	return model, nil
}

func fetchConnectionsGroups(date string) ([4]WordGroup, error) {
	url := fmt.Sprintf("https://www.nytimes.com/svc/connections/v2/%s.json", date)
	var groups [4]WordGroup

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		return groups, fmt.Errorf("error fetching Connections game: %v", err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return groups, fmt.Errorf("non-OK HTTP status: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return groups, fmt.Errorf("error reading response body: %v", err)
	}

	// Decode JSON
	var result ConnectionsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return groups, fmt.Errorf("error decoding JSON: %v\nbody: %s", err, string(body))
	}

	for i, category := range result.Categories {
		var members [4]string
		for i, card := range category.Cards {
			members[i] = card.Content
		}

		groups[i] = WordGroup{
			Members:    members,
			Clue:       category.Title,
			Color:      i + 1,
			IsRevealed: false,
		}
	}

	return groups, nil
}

// SaveToFile persists the current connections game state to a SQLite database.
func (m *ConnectionsModel) SaveToFile() error {
	// Create the data directory if it doesn't exist
	os.MkdirAll("data/connections", 0755)

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
