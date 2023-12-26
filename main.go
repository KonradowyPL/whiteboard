package main

import (
	"image/color"
	"log"
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

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// converts tile index to cords within chunk
func tileToCords(index byte) (byte, byte) {
	x := (index & 0x0F)
	y := (index & 0xF0) >> 4
	return x, y
}

func abs[T Number](x T) T {
	if x < T(0) {
		return -x
	}
	return x
}

// converts cords within chunk into tile index
//
// WARNING: does not check if given cords are within chunk
//
// inverse of func tileToCords()
func cordsToTile(x int, y int) byte {
	return byte((x & 0x0F) + (y&0x0F)<<4)
}

// converts global cords to chunk cords
func cordsToChunk(x int, y int) (int, int) {
	return x >> 4, y >> 4
}

// converts chunk cordinates into chunk id in memory
func (g *Game) getChunkAt(x int, y int) int {
	for i, chunk := range g.world.chunks {
		if chunk.x == x && chunk.y == y {
			return i
		}
	}
	return -1
}

// converts world cordinates into in chunk cordinates and to chunk cordinates
func cordsToPos(x int, y int) (inChunkX int, inChunkY int, chunkX int, chunkY int) {
	chunkX, chunkY = cordsToChunk(x, y)
	inChunkX = x & 0xF
	inChunkY = y & 0xF

	return inChunkX, inChunkY, chunkX, chunkY
}

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
	g.Editor(cursorPos)

	g.lastMouse = cursorPos
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

var white *ebiten.Image = ebiten.NewImage(32, 32)
var BuildTime string
var Version string

func main() {
	log.Println("Build time:", BuildTime)
	log.Println("Version:", Version)

	white.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})
	g := &Game{}
	g.camera.zoom = 1

	// temp code for generating chunks
	g.world.chunks = append(g.world.chunks, tempCreateChunk(-1, 0))
	g.world.chunks = append(g.world.chunks, tempCreateChunk(0, 0))
	g.world.chunks = append(g.world.chunks, tempCreateChunk(1, 0))
	g.world.chunks = append(g.world.chunks, tempCreateChunk(0, 1))
	g.world.chunks = append(g.world.chunks, tempCreateChunk(1, 1))

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// we can do this becouse in render we are filling screen it already
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowTitle("Hello, World!")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func tempCreateChunk(x, y int) loadedChunk {
	chunk := loadedChunk{x: x, y: y}

	for i := 0; i < len(chunk.grid); i++ {
		chunk.grid[i].Type = 0
	}
	return chunk
}
