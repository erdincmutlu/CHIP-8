package main

import (
	"fmt"
	"log"
	"os"

	"github.com/erdincmutlu/CHIP-8/c8"
	"github.com/hajimehoshi/ebiten"
)

var prog *c8.Prog

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage \"go run main.go ROM_NAME\"\n")
		return
	}

	err := c8.ReadROM(os.Args[1])
	if err != nil {
		return
	}

	go c8.RunROM()

	prog, err = c8.NewProg()
	if err != nil {
		log.Fatal(err)
	}

	err = ebiten.Run(update, c8.BoardWidth, c8.BoardHeight, 1, "Chip 8 - "+os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	if err := prog.Update(); err != nil {
		return err
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	prog.Draw(screen)
	return nil
}
