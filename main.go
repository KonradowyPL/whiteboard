package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"mrogalski.eu/go/vec"
)

type Game struct {
	camera    Camera
	lastMouse vec.Vec
}

type object struct {
	color color.RGBA
}

type Camera struct {
	pos  vec.Vec
	zoom float64
}

var grid [][]object

func (g *Game) Update() error {

	x, y := ebiten.CursorPosition()
	cursorPos := vec.New(float64(x), float64(y))

	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		delta := cursorPos.Sub(g.lastMouse)

		_, dy := ebiten.Wheel()

		// avoid division by 0
		if dy != 0 {

			// check posiiton mouse is hovering at before change
			before := cursorPos.Scale(1 / g.camera.zoom).Add(g.camera.pos)

			g.camera.zoom += dy
			if g.camera.zoom <= 0 {
				g.camera.zoom = 1
			}

			// check posiiton mouse is hovering at after change
			after := cursorPos.Scale(1 / g.camera.zoom).Add(g.camera.pos)

			// calculate diffrence and add it to camera pos
			// this is to coursor stays in the same place after zooming
			diff := before.Sub(after)
			g.camera.pos = g.camera.pos.Add(diff)
		}

		g.camera.pos = g.camera.pos.Add(delta.Scale(-1 / g.camera.zoom))
	}

	g.lastMouse = cursorPos

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	for x, row := range grid {
		for y, tile := range row {
			square := ebiten.NewImage(32, 32)
			square.Fill(tile.color)

			opts := &ebiten.DrawImageOptions{}

			opts.GeoM.Translate(float64(x)*32-g.camera.pos.X, float64(y)*32-g.camera.pos.Y)

			opts.GeoM.Scale(g.camera.zoom, g.camera.zoom)

			screen.DrawImage(square, opts)

		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%dTPS\n%.2fFPS", ebiten.TPS(), ebiten.ActualFPS()))

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func main() {

	for i := 0; i < 10; i++ {
		row := make([]object, 10)
		for j := 0; j < 10; j++ {
			r, g, b := uint8(rand.Intn(0xff)), uint8(rand.Intn(0xff)), uint8(rand.Intn(0xff))
			row[j] = object{color.RGBA{r, g, b, 0}}
		}
		grid = append(grid, row)
	}

	ebiten.SetTPS(20)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Hello, World!")

	if err := ebiten.RunGame(&Game{Camera{vec.New(0, 0), 1}, vec.New(0, 0)}); err != nil {
		log.Fatal(err)
	}
}
