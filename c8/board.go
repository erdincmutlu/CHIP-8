package c8

const (
	pixelSize          = 10
	pixelsHorizontally = 64
	pixelsVertically   = 32

	BoardWidth  = pixelsHorizontally * pixelSize
	BoardHeight = pixelsVertically * pixelSize
)

var Tiles [pixelsVertically][pixelsHorizontally]byte

type board struct {
	tiles [pixelsVertically][pixelsHorizontally]byte
}

var b2 board

// newBoard is to create board
func newBoard() (*board, error) {
	return &b2, nil
}
