package gogame

import (
	"slices"
)

// Board sizes
const (
	BoardSize9  = 9
	BoardSize13 = 13
	BoardSize19 = 19
	DefaultSize = BoardSize9
)

// Stone colors
const (
	Black = -1
	Empty = 0
	White = 1
)

// Group represents a connected set of stones of the same color.
type Group struct {
	Stones    []Position
	Liberties int
	Color     int8
}

// Board stores the game state.
type Board struct {
	Size  int
	Cells [][]int8
}

type Position struct {
	X, Y int
}

// NewBoard creates a new empty board of the given size.
func NewBoard(size int) Board {
	cells := make([][]int8, size)
	for y := range cells {
		cells[y] = make([]int8, size)
	}
	return Board{Size: size, Cells: cells}
}

// Clone creates a deep copy of the board.
func (b Board) Clone() Board {
	cells := make([][]int8, b.Size)
	for y := range cells {
		cells[y] = make([]int8, b.Size)
		copy(cells[y], b.Cells[y])
	}
	return Board{Size: b.Size, Cells: cells}
}

// Equal compares two boards for equality.
func (b Board) Equal(other Board) bool {
	for y := 0; y < b.Size; y++ {
		for x := 0; x < b.Size; x++ {
			if b.Cells[y][x] != other.Cells[y][x] {
				return false
			}
		}
	}
	return true
}

// HasStones checks if the board has any stones.
func (b Board) HasStones() bool {
	for y := 0; y < b.Size; y++ {
		for x := 0; x < b.Size; x++ {
			if b.Cells[y][x] != Empty {
				return true
			}
		}
	}
	return false
}

// InBounds checks if a position is within the board.
func (b Board) InBounds(pos Position) bool {
	return pos.X >= 0 && pos.X < b.Size && pos.Y >= 0 && pos.Y < b.Size
}

// PlaceStone places a stone at the given position and returns the captured positions.
func (b *Board) PlaceStone(pos Position, color int8) []Position {
	b.Cells[pos.Y][pos.X] = color

	// Check for captures of opponent groups adjacent to this placement
	opponent := color * -1
	var captured []Position
	adjacent := []Position{
		{pos.X - 1, pos.Y},
		{pos.X + 1, pos.Y},
		{pos.X, pos.Y - 1},
		{pos.X, pos.Y + 1},
	}

	for _, adj := range adjacent {
		// Skip if adjacent position is out of bounds or not opponent
		if !b.InBounds(adj) || b.Cells[adj.Y][adj.X] != opponent {
			continue
		}

		// Skip if we already captured this group
		if slices.Contains(captured, adj) {
			continue
		}

		// Capture group if it has 0 liberties
		group := b.GetGroup(adj)
		if group.Liberties == 0 {
			b.RemoveGroup(group)
			captured = append(captured, group.Stones...)
		}
	}

	return captured
}

// GetGroup performs BFS to find all connected stones and count the groups liberties.
func (b Board) GetGroup(pos Position) Group {
	if !b.InBounds(pos) || b.Cells[pos.Y][pos.X] == Empty {
		return Group{}
	}

	color := b.Cells[pos.Y][pos.X]
	visited := make(map[Position]bool)
	var stones []Position
	liberties := 0
	libertySet := make(map[Position]bool)

	queue := []Position{pos}
	for len(queue) > 0 {
		// Pop next position from queue
		p := queue[0]
		queue = queue[1:]

		if visited[p] {
			continue
		}

		visited[p] = true

		if b.Cells[p.Y][p.X] == color {
			stones = append(stones, p)

			// Check adjacent positions
			adjacent := []Position{
				{p.X - 1, p.Y},
				{p.X + 1, p.Y},
				{p.X, p.Y - 1},
				{p.X, p.Y + 1},
			}
			for _, adj := range adjacent {
				if !b.InBounds(adj) || visited[adj] {
					continue
				}

				adjacentPos := b.Cells[adj.Y][adj.X]

				// If adjacent position is empty, add it as a liberty
				if adjacentPos == Empty {
					if !libertySet[adj] {
						libertySet[adj] = true
						liberties++
					}
				}

				// If adjacent position is the same color, add it to queue
				if adjacentPos == color {
					queue = append(queue, adj)
				}
			}
		}
	}

	return Group{
		Stones:    stones,
		Liberties: liberties,
		Color:     color,
	}
}

// RemoveGroup removes all stones in the group from the board.
func (b *Board) RemoveGroup(group Group) {
	for _, pos := range group.Stones {
		b.Cells[pos.Y][pos.X] = Empty
	}
}
