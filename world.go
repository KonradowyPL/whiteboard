package main

import (
	"image/color"

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

	square Square // 1
	circle Circle // 2
}

type Square struct {
	color color.RGBA
}

func (t *Square) render(g *Game, screen *ebiten.Image, x float64, y float64) {
	img := ebiten.NewImage(32, 32)
	img.Fill(t.color)

	g.basicRedner(screen, x, y, img)
}

type Circle struct {
	color color.RGBA
}

func (t *Circle) render(g *Game, screen *ebiten.Image, x float64, y float64) {
	img := ebiten.NewImage(32, 32)
	img.Fill(t.color)

	g.basicRedner(screen, x, y, img)
}

func (w *World) Update() error {
	return nil
}
