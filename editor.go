package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"mrogalski.eu/go/vec"
)

var lastPLaceX int
var lastPLaceY int

func (g *Game) Editor(cursor vec.Vec) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		world := g.screenToWorldspace(cursor).Scale(0.03125)
		x, y := int(world.X), int(world.Y)

		g.placeBeltAt(x, y)

	}
}

func (g *Game) placeBeltAt(targetX, targetY int) {

	var x, y int = targetX, targetY
	var dirty = false

	if targetX > (lastPLaceX + 7) {
		x = lastPLaceX + 7
		dirty = true
	}

	if targetY > (lastPLaceY + 7) {
		y = lastPLaceY + 7
		dirty = true
	}

	if targetX < (lastPLaceX - 7) {
		x = lastPLaceX - 7
		dirty = true
	}

	if targetY < (lastPLaceY - 7) {
		y = lastPLaceY - 7
		dirty = true
	}

	inChunkX, inChunkY, chunkX, chunkY := cordsToPos(x, y)

	tileIdx := cordsToTile(inChunkX, inChunkY)

	chunkIdx := g.getChunkAt(chunkX, chunkY)

	g.world.chunks[chunkIdx].grid[tileIdx].Type = 1
	g.world.chunks[chunkIdx].grid[tileIdx].belt.next = tileIdx

	lastPLaceChunkX, lastPLaceChunkY := cordsToChunk(lastPLaceX, lastPLaceY)
	lastPLaceChunkIdx := g.getChunkAt(lastPLaceChunkX, lastPLaceChunkY)
	lastPlace := cordsToTile(lastPLaceX, lastPLaceY)

	g.world.chunks[lastPLaceChunkIdx].grid[lastPlace].belt.next = tileIdx

	lastPLaceX = x
	lastPLaceY = y

	if dirty {
		g.placeBeltAt(targetX, targetY)
	}
}
