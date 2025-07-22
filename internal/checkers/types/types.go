package types

type Position struct {
	X, Y int
}

type Piece struct {
	Color int8
	Value int8
}

type CaptureMove struct {
	From, To Position
}

type RenderContext struct {
	Board        [][]Piece
	Selected     Position
	ValidMoves   []Position
	CaptureMoves []CaptureMove
}
