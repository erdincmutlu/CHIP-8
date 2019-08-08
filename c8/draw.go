package c8

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

const ()

var (
	backgroundColor = color.Black
)

var tiles [tilesVertically][tilesHorizontally]byte

var tile *ebiten.Image

// Prog represent a program state
type Prog struct {
	board     *board
	progImage *ebiten.Image
}

// NewProg generates a new Prog object
func NewProg() (*Prog, error) {
	progImage, _ := ebiten.NewImage(10, 10, ebiten.FilterDefault)
	theBoard, _ := newBoard()
	p := &Prog{
		progImage: progImage,
		board:     theBoard,
	}

	tile, _ = ebiten.NewImage(10, 10, ebiten.FilterDefault)
	return p, nil
}

// Update is to update screen
func (p *Prog) Update() error {
	//ebitenutil.DebugPrint(screen, "Hello World!")
	// screen.Fill(color.Black)
	return nil
}

// Draw draws the current game to the given screen
func (p *Prog) Draw(screen *ebiten.Image) error {
	screen.Fill(backgroundColor)

	for row := 0; row < tilesVertically; row++ {
		for col := 0; col < tilesHorizontally; col++ {
			digits := getDigits(tiles[row][col])
			for i, digit := range digits {
				if digit {
					drawATile(screen, row, col, i)
				}
			}
		}
	}
	drawATile(screen, 0, 0, 0)
	drawATile(screen, 31, 7, 7)
	return nil
}

func drawATile(screen *ebiten.Image, row, col, index int) error {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64((col*8+index)*tileSize), float64(row*tileSize))
	op.ColorM.Translate(0xFF, 0x10, 0x20, 0xBB)
	screen.DrawImage(tile, op)
	return nil
}

func drawScreen() {
	// time.Sleep(time.Second * 5)
	// fmt.Printf(clearScreen)
	fmt.Printf("+----------------------------------------------------------------+\n")
	for row := 0; row < tilesVertically; row++ {
		fmt.Printf("|")
		for col := 0; col < tilesHorizontally; col++ {
			fmt.Printf("%s", getByteForScreen(tiles[row][col]))
		}
		fmt.Printf("|\n")
	}
	fmt.Printf("+----------------------------------------------------------------+\n")
}

func getDigits(x byte) [8]bool {
	var val [8]bool
	index := 7
	for x > 0 {
		if x%2 == 1 {
			val[index] = true
		}
		x = x >> 1
		index--
	}
	return val
}
