package c8

const (
	tileSize          = 10
	tilesHorizontally = 64 / 8
	tilesVertically   = 32

	BoardWidth  = tilesHorizontally * 8 * tileSize
	BoardHeight = tilesVertically * tileSize
)

var Tiles [tilesVertically][tilesHorizontally]byte

type board struct {
	tiles [tilesVertically][tilesHorizontally]byte
}

var b2 board

// newBoard is to create board
func newBoard() (*board, error) {
	return &b2, nil
}
