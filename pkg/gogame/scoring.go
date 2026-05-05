package gogame

// Komi is the compensation points given to White for going second.
const Komi = 7.5

// Score holds the result of a completed game.
type Score struct {
	BlackScore float64
	WhiteScore float64
	Winner     int8
}

// CalculateScore computes the area score for both players.
// Area scoring = territory + stones on board.
func CalculateScore(board Board) Score {
	territoryBlack, territoryWhite := countTerritory(board)

	// Count stones on the board
	stonesBlack := 0
	stonesWhite := 0
	for y := 0; y < board.Size; y++ {
		for x := 0; x < board.Size; x++ {
			switch board.Cells[y][x] {
			case Black:
				stonesBlack++
			case White:
				stonesWhite++
			}
		}
	}

	blackScore := float64(territoryBlack + stonesBlack)
	whiteScore := float64(territoryWhite+stonesWhite) + Komi

	var winner int8 = Empty
	if blackScore > whiteScore {
		winner = Black
	} else if whiteScore > blackScore {
		winner = White
	}

	return Score{
		BlackScore: blackScore,
		WhiteScore: whiteScore,
		Winner:     winner,
	}
}

// countTerritory performs flood fill to determine territory ownership.
func countTerritory(board Board) (int, int) {
	visited := make(map[Position]bool)
	blackTerritory := 0
	whiteTerritory := 0

	for y := 0; y < board.Size; y++ {
		for x := 0; x < board.Size; x++ {
			pos := Position{X: x, Y: y}
			if visited[pos] {
				continue
			}
			if board.Cells[y][x] != Empty {
				continue
			}

			// Flood fill this empty region
			region := make([]Position, 0)
			touchesBlack := false
			touchesWhite := false
			regionVisited := make(map[Position]bool)

			queue := []Position{pos}
			for len(queue) > 0 {
				// Pop next position from queue
				p := queue[0]
				queue = queue[1:]

				// Check stone color
				switch board.Cells[p.Y][p.X] {
				case Black:
					touchesBlack = true
					continue
				case White:
					touchesWhite = true
					continue
				}

				// Check if this position has already been visited
				if regionVisited[p] || visited[p] {
					continue
				}
				regionVisited[p] = true
				region = append(region, p)

				// Add adjacent positions to the queue
				adjacent := []Position{
					{p.X - 1, p.Y},
					{p.X + 1, p.Y},
					{p.X, p.Y - 1},
					{p.X, p.Y + 1},
				}
				for _, adj := range adjacent {
					if board.InBounds(adj) && !regionVisited[adj] {
						queue = append(queue, adj)
					}
				}
			}

			// Mark all positions in this region as visited
			for _, p := range region {
				visited[p] = true
			}

			// Determine territory ownership
			if touchesBlack && !touchesWhite {
				blackTerritory += len(region)
			} else if touchesWhite && !touchesBlack {
				whiteTerritory += len(region)
			}
			// If both or neither, it's neutral (dame)
		}
	}

	return blackTerritory, whiteTerritory
}
