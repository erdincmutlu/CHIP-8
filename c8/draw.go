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

var pixels [pixelsVertically][pixelsHorizontally]byte

var aPixel *ebiten.Image

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

	aPixel, _ = ebiten.NewImage(10, 10, ebiten.FilterDefault)
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

	for row := 0; row < pixelsVertically; row++ {
		for col := 0; col < pixelsHorizontally; col++ {
			if pixels[row][col] > 0 {
				drawAPixel(screen, row, col)
			}
		}
	}
	// drawAPixel(screen, 0, 0)
	// drawAPixel(screen, 1, 1)
	// drawAPixel(screen, 31, 63)
	return nil
}

func drawAPixel(screen *ebiten.Image, row, col int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(col*pixelSize), float64(row*pixelSize))
	op.ColorM.Translate(0xFF, 0x10, 0x20, 0xBB)
	screen.DrawImage(aPixel, op)
}

func drawScreen() {
	// time.Sleep(time.Second * 5)
	// fmt.Printf(clearScreen)
	fmt.Printf("+----------------------------------------------------------------+\n")
	for row := 0; row < pixelsVertically; row++ {
		for col := 0; col < pixelsHorizontally; col++ {
			if pixels[row][col] > 0 {
				fmt.Printf("%d", pixels[row][col])
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("+----------------------------------------------------------------+\n")
}
