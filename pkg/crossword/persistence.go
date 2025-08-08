package crossword

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

const filename = "data/crossword/games.db"

// JSON response objects from NYT
type PuzzleResponse struct {
	Body []PuzzleBody `json:"body"`
}

type PuzzleBody struct {
	Cells      []Cell     `json:"cells"`
	Clues      []Clue     `json:"clues"`
	Dimensions Dimensions `json:"dimensions"`
}

type Cell struct {
	Answer string `json:"answer"`
	Label  string `json:"label"`
	Type   int    `json:"type"`
	Clues  []int  `json:"clues"`
}

type Clue struct {
	Label     string     `json:"label"`
	Direction string     `json:"direction"`
	Text      []ClueText `json:"text"`
}

type ClueText struct {
	Plain string `json:"plain"`
}

type Dimensions struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

// LoadGame returns the Wordle game state for a given date.
func LoadGame(date string) (CrosswordModel, error) {
	// Try loading the saved game from the database
	model, err := LoadFromFile(date)
	if err == nil {
		return model, nil
	}

	// Try fetching from the web if not in database
	model, fetchErr := fetchCrosswordGame(date)
	if fetchErr != nil {
		return model, fetchErr
	}

	return model, nil
}

func fetchCrosswordGame(date string) (CrosswordModel, error) {
	var model CrosswordModel
	url := fmt.Sprintf("https://www.nytimes.com/svc/crosswords/v6/puzzle/daily/%s.json", date)

	// Make the GET request
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Games-Auth-Bypass", "true")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return model, fmt.Errorf("error fetching Connections game: %v", err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return model, fmt.Errorf("non-OK HTTP status: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model, fmt.Errorf("error reading response body: %v", err)
	}

	// Decode JSON
	var result PuzzleResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return model, fmt.Errorf("error decoding JSON: %v\nbody: %s", err, string(body))
	}
	data := result.Body[0]

	height := data.Dimensions.Height
	width := data.Dimensions.Width

	// Extract Across and Down clues
	var acrossClues, downClues []string
	for _, clue := range data.Clues {
		// Join all text segments
		full := ""
		for _, seg := range clue.Text {
			full += seg.Plain
		}

		entry := fmt.Sprintf("%s %s", clue.Label, full)

		if clue.Direction == "Across" {
			acrossClues = append(acrossClues, entry)
		} else {
			downClues = append(downClues, entry)
		}
	}

	// Initialize grid structures
	answer := make([][]byte, height)
	gridNums := make([][]int, height)
	clueIndices := make([][]Position, height)

	for i := range height {
		answer[i] = make([]byte, width)
		gridNums[i] = make([]int, width)
		clueIndices[i] = make([]Position, width)
	}

	// Parse grid data from response
	for idx, cell := range data.Cells {
		row := idx / width
		col := idx % width

		if len(cell.Answer) == 1 {
			answer[row][col] = cell.Answer[0]

			if len(cell.Label) > 0 {
				gridNums[row][col], _ = strconv.Atoi(cell.Label)
			}

			if len(cell.Clues) == 2 {
				clueIndices[row][col] = Position{
					X: cell.Clues[0],
					Y: cell.Clues[1] - len(acrossClues),
				}
			} else {
				clueIndices[row][col] = Position{
					X: cell.Clues[0],
					Y: -1,
				}
			}
		} else {
			answer[row][col] = '.'
			clueIndices[row][col] = Position{-1, -1}
		}
	}

	model.date = date
	model.acrossClues = acrossClues
	model.downClues = downClues
	model.answer = answer
	model.width = width
	model.height = height
	model.gridNums = gridNums
	model.clueIndices = clueIndices
	model.isAcross = true
	model.message = ""

	model.handleReset()
	model.prepareGrid()

	return model, nil
}

// SaveToFile persists the current crossword puzzle state to a SQLite database.
func (m *CrosswordModel) SaveToFile() error {
	// Create the data directory if it doesn't exist
	os.MkdirAll("data/crossword", 0755)

	// Open the database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS crosswords (
		date TEXT PRIMARY KEY,
		across TEXT,
		down TEXT,
		answer TEXT,
		grid TEXT,
		grid_nums TEXT,
		clue_indices TEXT,
		width INT,
		height INT
	)`
	if _, err := db.Exec(createTableQuery); err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	// Convert slices to JSON
	acrossJSON, _ := json.Marshal(m.acrossClues)
	downJSON, _ := json.Marshal(m.downClues)
	answerJSON, _ := json.Marshal(m.answer)
	gridJSON, _ := json.Marshal(m.grid)
	gridNumsJSON, _ := json.Marshal(m.gridNums)
	clueIndicesJSON, _ := json.Marshal(m.clueIndices)
	widthJSON, _ := json.Marshal(m.width)
	heightJSON, _ := json.Marshal(m.height)

	// Insert the data into the database
	insertQuery := `
		INSERT OR REPLACE INTO crosswords(
			date, across, down, answer, grid, grid_nums, clue_indices, width, height
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(insertQuery,
		m.date,
		acrossJSON,
		downJSON,
		answerJSON,
		gridJSON,
		gridNumsJSON,
		clueIndicesJSON,
		widthJSON,
		heightJSON,
	)

	return err
}

// LoadFromFile loads a crossword puzzle state from the SQLite database.
func LoadFromFile(date string) (CrosswordModel, error) {
	var model CrosswordModel

	// Open the database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return model, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	query := `
		SELECT across, down, answer, grid, grid_nums, clue_indices, width, height
		FROM crosswords
		WHERE date = ?
	`

	var acrossJSON, downJSON []byte
	var answerJSON, gridJSON []byte
	var gridNumsJSON, clueIndicesJSON []byte
	var widthJSON, heightJSON []byte

	// Query row from database
	err = db.QueryRow(query, date).Scan(
		&acrossJSON, &downJSON,
		&answerJSON, &gridJSON,
		&gridNumsJSON, &clueIndicesJSON,
		&widthJSON, &heightJSON,
	)
	if err != nil {
		return model, fmt.Errorf("failed to fetch saved puzzle: %v", err)
	}

	// Convert JSON to slices
	var across, down []string
	var answer, grid [][]byte
	var gridNums [][]int
	var clueIndices [][]Position
	var width, height int

	json.Unmarshal(acrossJSON, &across)
	json.Unmarshal(downJSON, &down)
	json.Unmarshal(answerJSON, &answer)
	json.Unmarshal(gridJSON, &grid)
	json.Unmarshal(gridNumsJSON, &gridNums)
	json.Unmarshal(clueIndicesJSON, &clueIndices)
	json.Unmarshal(widthJSON, &width)
	json.Unmarshal(heightJSON, &height)

	// Set up the model with data from the saved game
	model.date = date
	model.acrossClues = across
	model.downClues = down
	model.answer = answer
	model.grid = grid
	model.gridNums = gridNums
	model.clueIndices = clueIndices
	model.width = width
	model.height = height
	model.isAcross = true
	model.message = ""

	// Initialize grid numbers, clue indices, and solved status
	model.prepareGrid()

	return model, nil
}

// prepareGrid initializes the crossword model's grid state,
// tracking solved clues, black squares, and progress metrics.
func (m *CrosswordModel) prepareGrid() {
	m.incorrect = make([][]bool, m.height)
	for i := range m.incorrect {
		m.incorrect[i] = make([]bool, m.width)
	}

	// Iterate through the grid row by row
	for row := range m.height {
		for col := range m.width {

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

			// Check if this cell starts an across clue and if the clue is solved
			if col == 0 || m.grid[row][col-1] == '.' {
				isSolved := true
				for x := col; x < m.width && m.grid[row][x] != '.'; x++ {
					if m.grid[row][x] == ' ' {
						isSolved = false
						break
					}
				}
				m.isAcrossSolved = append(m.isAcrossSolved, isSolved)
			}

			// Check if this cell starts a down clue and if the clue is solved
			if row == 0 || m.grid[row-1][col] == '.' {
				isSolved := true
				for y := row; y < m.height && m.grid[y][col] != '.'; y++ {
					if m.grid[y][col] == ' ' {
						isSolved = false
					}
				}
				m.isDownSolved = append(m.isDownSolved, isSolved)
			}
		}
	}
}
