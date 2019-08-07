package c8

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	BoardWidth = 420
	BoardHeight = 600
)

func Update(screen *ebiten.Image) error {
	ebitenutil.DebugPrint(screen, "Hello World!")
	return nil
}
