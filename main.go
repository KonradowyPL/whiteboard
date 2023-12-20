package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"mrogalski.eu/go/vec"
)

// TODO
// low: clean up variable declatarions
type Game struct {
	camera    Camera
	lastMouse vec.Vec

	world World
}

type Camera struct {
	pos  vec.Vec
	zoom float64
}

var zoom = struct {
	min float64
	max float64
}{0.05, 4}

// converts screen position in pixels to world space
func (g *Game) screenToWorldspace(pos vec.Vec) vec.Vec {
	return pos.Scale(1 / g.camera.zoom).Add(g.camera.pos)
}

// converts world cordinates to screen position in pixels
func (g *Game) worldToScreen(pos vec.Vec) vec.Vec {
	return g.worldToGlobal(pos).Scale(g.camera.zoom)
}

// makes cordinates relateive to camerea
func (g *Game) worldToGlobal(pos vec.Vec) vec.Vec {
	return pos.Sub(g.camera.pos)
}

func (g *Game) moveCamera(cursorPos vec.Vec) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		delta := cursorPos.Sub(g.lastMouse)

		_, dy := ebiten.Wheel()

		// avoid division by 0
		if dy != 0 {

			// check posiiton mouse is hovering at before change
			before := g.screenToWorldspace(cursorPos)

			g.camera.zoom *= dy*.1 + 1
			// check if within bounds
			if g.camera.zoom <= zoom.min {
				g.camera.zoom = zoom.min
			} else if g.camera.zoom >= zoom.max {
				g.camera.zoom = zoom.max
			}

			// check posiiton mouse is hovering at after change
			after := g.screenToWorldspace(cursorPos)

			// calculate diffrence and add it to camera pos
			// this is to coursor stays in the same place after zooming
			diff := before.Sub(after)
			g.camera.pos = g.camera.pos.Add(diff)

			// TODO
			// low: fixing camera alignment proces proably can be optymalised
		}

		g.camera.pos = g.camera.pos.Add(delta.Scale(-1 / g.camera.zoom))
	}
}

func (g *Game) Update() error {
	now := time.Now().UnixNano()

	// delta time in seconds
	// 1e-9 = nano prefix
	deltaTime := float64((now - g.world.lastUpdate)) * 1e-9

	// 20Hz
	if deltaTime >= 0.05 {
		g.world.lastUpdate = now
		// update world
		// world is ran on 20Hz while
		// while game is rendered at
		err := g.world.Update()
		if err != nil {
			return err
		}
	}

	x, y := ebiten.CursorPosition()
	cursorPos := vec.New(float64(x), float64(y))
	g.moveCamera(cursorPos)

	g.lastMouse = cursorPos
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func main() {
	game := &Game{}
	game.camera.zoom = 1

	for i := 0; i < 10; i++ {
		row := make([]object, 10)
		for j := 0; j < 10; j++ {
			r, g, b := uint8(rand.Intn(0xff)), uint8(rand.Intn(0xff)), uint8(rand.Intn(0xff))
			row[j] = object{color.RGBA{r, g, b, 0xff}}
		}
		game.world.grid = append(game.world.grid, row)
	}

	ebiten.SetTPS(60)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Hello, World!")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
