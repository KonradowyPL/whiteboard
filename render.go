package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"mrogalski.eu/go/vec"
)

// debug funciton to print 8x8 pixel on the screen
func (g *Game) debugRender(screen *ebiten.Image, pos vec.Vec) {
	square := ebiten.NewImage(8, 8)
	square.Fill(color.White)

	opts := &ebiten.DrawImageOptions{}

	screenPos := g.worldToGlobal(pos)

	opts.GeoM.Translate(screenPos.X, screenPos.Y)
	opts.GeoM.Scale(g.camera.zoom, g.camera.zoom)

	screen.DrawImage(square, opts)

}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(HexToRGBA(0x010730ff))
	g.renderGrid(screen)

	for i, _ := range g.world.chunks {
		g.renderChunk(screen, i)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%dTPS\n%.2fFPS", ebiten.TPS(), ebiten.ActualFPS()))
}

// renders single chunk
func (g *Game) renderChunk(screen *ebiten.Image, id int) {
	chunk := g.world.chunks[id]
	x := float64(chunk.x << 9)
	y := float64(chunk.y << 9)

	for i, tile := range chunk.grid {
		tileX, tileY := tileToCords(byte(i))

		square := ebiten.NewImage(32, 32)
		square.Fill(tile.color)

		opts := &ebiten.DrawImageOptions{}

		screenPos := g.worldToGlobal(vec.New(float64(tileX)*32+x, float64(tileY)*32+y))

		opts.GeoM.Translate(screenPos.X, screenPos.Y)
		opts.GeoM.Scale(g.camera.zoom, g.camera.zoom)

		screen.DrawImage(square, opts)

	}
}

func (g *Game) renderGrid(screen *ebiten.Image) {
	inv := int32(max(1/g.camera.zoom, 1) * 255)

	// draw grid
	// grid is displayed as other grids containing 8x8 slice of previous grid
	g.drawLines(screen, 32, HexToRGBA(0x0f153dff), uint32(max(512-inv, 0)))
	g.drawLines(screen, 256, HexToRGBA(0x1a2254ff), uint32(max(4096-inv, 0)>>3))
	g.drawLines(screen, 2048, HexToRGBA(0x232c68ff), uint32(max(32768-inv, 0)>>6))

}

func (g *Game) drawLines(screen *ebiten.Image, step float64, c color.Color, opacity uint32) {
	if opacity <= 0 {
		return
	}

	c = Blend(c, opacity)

	// get screen size
	size := screen.Bounds()
	sizeX := size.Dx()
	sizeY := size.Dy()

	// calculate starting position
	from := g.screenToWorldspace(vec.Zero).Scale(1 / step)
	to := g.screenToWorldspace(vec.New(float64(sizeX), float64(sizeY))).Scale(1 / step)

	// calculate
	fromX := math.Floor(from.X)
	fromY, _ := math.Modf(from.Y)

	// print vertical lines
	for x := fromX; x < to.X; x++ {
		xPos := g.worldToScreen(vec.New(x*step, 0)).X
		ebitenutil.DrawLine(screen, xPos, 0, xPos, float64(sizeY), c)
	}
	// print horizontal lines
	for y := fromY; y < to.Y; y++ {
		yPos := g.worldToScreen(vec.New(0, y*step)).Y
		ebitenutil.DrawLine(screen, 0, yPos, float64(sizeX), yPos, c)
	}
	// TODO
	// low: use somehting diffrent that 'ebitenutil.DrawLine'
}

// multiples given color by blend factor
func Blend(c color.Color, factor uint32) color.Color {
	factor = max(min(factor, 255), 0)
	r, g, b, a := c.RGBA()
	r = r * factor / 0xffff
	g = g * factor / 0xffff
	b = b * factor / 0xffff
	a = a * factor / 0xffff
	return color.RGBA{
		uint8(r),
		uint8(g),
		uint8(b),
		uint8(a),
	}
}

// converts hex to color.RGBA
func HexToRGBA(c uint32) color.RGBA {
	r := (c >> 24) & 0xFF
	g := (c >> 16) & 0xFF
	b := (c >> 8) & 0xFF
	a := c & 0xFF

	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}
