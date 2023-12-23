package main

import (
	"github.com/hajimehoshi/ebiten/v2"
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

	// 0 is reservated for air/null tiles

	belt Belt
}

type Belt struct {
	next byte
}

func (t *Belt) render(g *Game, screen *ebiten.Image, x float64, y float64) {
	g.basicRedner(screen, x, y, white)
}

func (w *World) Update() error {
	return nil
}
