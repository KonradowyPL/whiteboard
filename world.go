package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"mrogalski.eu/go/vec"
)

type World struct {
	lastUpdate int64

	chunks []loadedChunk
}

type loadedChunk struct {
	x int
	y int

	grid [256]object
}

type object struct {
	Type int

	// TYpe = 0 is reservated for air/null tiles

	belt Belt // Type = 1
}

type Belt struct {
	next byte
}

func (t *Belt) render(g *Game, screen *ebiten.Image, x float64, chunkX float64, y float64, chunkY float64) {
	g.basicRedner(screen, chunkX*512+x*32, chunkY*512+y*32, white)

	x2, y2 := tileToCords(t.next)

	var convertedX float64
	var convertedY float64

	// multi chunk owning
	// this code ensures that if one belt is on the edge of a chunk
	// it will render properly next child
	if x-float64(x2) >= 8 {
		convertedX = float64(x2) + 16
	} else if x-float64(x2) <= -8 {
		convertedX = float64(x2) - 16
	} else {
		convertedX = float64(x2)
	}

	if y-float64(y2) >= 8 {
		convertedY = float64(y2) + 16
	} else if y-float64(y2) <= -8 {
		convertedY = float64(y2) - 16
	} else {
		convertedY = float64(y2)
	}

	p1 := g.worldToScreen(vec.New(chunkX*512+float64(x)*32, chunkY*512+float64(y)*32))
	p2 := g.worldToScreen(vec.New(chunkX*512+convertedX*32, chunkY*512+convertedY*32))

	ebitenutil.DrawLine(screen, p1.X, p1.Y, p2.X, p2.Y, color.White)

}

func (w *World) Update() error {
	return nil
}
