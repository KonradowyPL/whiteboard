package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"mrogalski.eu/go/vec"
)

func (g *Game) Draw(screen *ebiten.Image) {

	for x, row := range g.world.grid {
		for y, tile := range row {
			square := ebiten.NewImage(32, 32)
			square.Fill(tile.color)

			opts := &ebiten.DrawImageOptions{}

			screenPos := g.worldToGlobal(vec.New(float64(x)*32, float64(y)*32))

			opts.GeoM.Translate(screenPos.X, screenPos.Y)
			opts.GeoM.Scale(g.camera.zoom, g.camera.zoom)

			screen.DrawImage(square, opts)

		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%dTPS\n%.2fFPS", ebiten.TPS(), ebiten.ActualFPS()))
}
