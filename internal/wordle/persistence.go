package wordle

import (
	"encoding/json"
	"fmt"
	"os"
)

const filename = "data/wordle/games.json"

// saveData represents the structure used to serialize Wordle data to JSON.
type saveData struct {
	Answer   string         `json:"answer"`
	Guesses  [6]string      `json:"guesses"`
	CursorX  int            `json:"cursor_x"`
	CursorY  int            `json:"cursor_y"`
	Keyboard map[string]int `json:"keyboard_state"`
}

// SaveToFile writes the current game state to a JSON file.
func (m *WordleModel) SaveToFile() error {
	// Read existing data
	savedGames := make(map[string]saveData)
	fileData, err := os.ReadFile(filename)
	if err == nil {
		if err := json.Unmarshal(fileData, &savedGames); err != nil {
			return fmt.Errorf("error parsing existing data: %v", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error reading file: %v", err)
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

	// Update or create entry for current date
	savedGames[m.date] = saveData{
		Answer:   string(m.answer[:]),
		Guesses:  guesses,
		CursorX:  m.cursorX,
		CursorY:  m.cursorY,
		Keyboard: keyboard,
	}

	// Write updated data back to file
	jsonData, err := json.MarshalIndent(savedGames, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}

	return os.WriteFile(filename, jsonData, 0644)
}

// LoadFromFile loads the most recent game state from the JSON file.
func LoadFromFile() (WordleModel, error) {
	// Read data
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return WordleModel{}, fmt.Errorf("error reading file: %w", err)
	}
	var data map[string]saveData
	if err := json.Unmarshal(fileContent, &data); err != nil {
		return WordleModel{}, fmt.Errorf("error parsing JSON: %w", err)
	}

	// Extract dates
	dates := make([]string, 0, len(data))
	for date := range data {
		dates = append(dates, date)
	}

	// Get most recent entry
	latestDate := dates[len(data)-1]
	entry := data[latestDate]

	// Convert to WordleModel
	model := WordleModel{
		date:     latestDate,
		answer:   [5]byte{},
		guesses:  [6][5]byte{},
		cursorX:  entry.CursorX,
		cursorY:  entry.CursorY,
		keyboard: make(map[byte]int, len(entry.Keyboard)),
		message:  "",
	}

	// Convert string data to byte arrays for answer and guesses
	copy(model.answer[:], entry.Answer)
	for i, guessStr := range entry.Guesses {
		copy(model.guesses[i][:], guessStr)
	}

	// Populate keyboard
	for key, value := range entry.Keyboard {
		model.keyboard[key[0]] = value
	}

	return model, nil
}
