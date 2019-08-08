package c8

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

const ()

var (
	backgroundColor = color.RGBA{0x0, 0x0, 0x0, 0xff}
)

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
	// screen.Fill(backgroundColor)
	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(float64(100), float64(150))
	// op.ColorM.Translate(0xFF, 0x10, 0x20, 0xBB)
	// screen.DrawImage(tile, op)
	// op.GeoM.Translate(float64(120), float64(170))
	// screen.DrawImage(tile, op)

	drawSingle(screen, 0, 0, 0)
	drawSingle(screen, 1, 1, 1)
	return nil
}

func drawSingle(screen *ebiten.Image, row, col, index int) error {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(row*tileSize), float64((col*4+index)*tileSize))
	op.ColorM.Translate(0xFF, 0x10, 0x20, 0xBB)
	screen.DrawImage(tile, op)
	return nil
}

func drawScreen() {
	// time.Sleep(time.Second * 2)
	// fmt.Printf(clearScreen)
	fmt.Printf("+----------------------------------------------------------------+\n")
	for row := 0; row < ScreenHeight; row++ {
		fmt.Printf("|")
		for col := 0; col < ScreenWidth; col++ {
			fmt.Printf("%s", getByteForScreen(screen[row][col]))
		}
		fmt.Printf("|\n")
	}
	fmt.Printf("+----------------------------------------------------------------+\n")
}
