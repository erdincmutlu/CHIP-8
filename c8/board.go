package c8

const (
	BoardWidth   = 600
	BoardHeight  = 420
	tileSize     = 10
	ScreenWidth  = 64 / 8
	ScreenHeight = 32
)

var screen [ScreenHeight][ScreenWidth]byte

type board struct {
	screen2 [ScreenHeight][ScreenWidth]byte
}

var b2 board

// newBoard is to create board
func newBoard() (*board, error) {
	return &b2, nil
}
