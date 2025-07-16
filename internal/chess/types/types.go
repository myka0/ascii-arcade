package types

type Position struct {
	X, Y int
}

type Piece struct {
	Color int8
	Value int8
}

type RenderContext struct {
	Board      [][]Piece
	Selected   Position
	ValidMoves []Position

	WhiteCapturedPieces []Piece
	BlackCapturedPieces []Piece

	IsWhiteKingInCheck bool
	IsBlackKingInCheck bool
}
