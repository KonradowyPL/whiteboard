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

func (t *Belt) render(g *Game, screen *ebiten.Image, x int, chunkX int, y int, chunkY int) {
	g.basicRedner(screen, float64(chunkX*512+x*32), float64(chunkY*512+y*32), white)

	x2, y2 := tileToCords(t.next)

	var convertedX int
	var convertedY int

	// multi chunk owning
	// this code ensures that if one belt is on the edge of a chunk
	// it will render properly next child
	if x-int(x2) >= 8 {
		convertedX = int(x2) + 16
	} else if x-int(x2) <= -8 {
		convertedX = int(x2) - 16
	} else {
		convertedX = int(x2)
	}

	if y-int(y2) >= 8 {
		convertedY = int(y2) + 16
	} else if y-int(y2) <= -8 {
		convertedY = int(y2) - 16
	} else {
		convertedY = int(y2)
	}

	for xPos := min(x, int(x2)); xPos < max(x, int(x2)); xPos++ {
		g.basicRedner(screen, float64(chunkX*512+xPos*32), float64(chunkY*512+y*32), white)

	}

	p1 := g.worldToScreen(vec.New(float64(chunkX*512+x*32), float64(chunkY*512+y*32)))
	p2 := g.worldToScreen(vec.New(float64(chunkX*512+convertedX*32), float64(chunkY*512+convertedY*32)))

	ebitenutil.DrawLine(screen, p1.X, p1.Y, p2.X, p2.Y, color.White)

}

func (w *World) Update() error {
	return nil
}
