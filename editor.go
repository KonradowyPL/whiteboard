package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"mrogalski.eu/go/vec"
)

func (g *Game) Editor(cursor vec.Vec) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		world := g.screenToWorldspace(cursor).Scale(0.03125)
		x, y := int(world.X), int(world.Y)

		chunkX, chunkY := cordsToChunk(x, y)
		inChunkX := x & 0xF
		inChunkY := y & 0xF

		tileIdx := cordsToTile(byte(inChunkX), byte(inChunkY))

		chunkIdx := g.getChunkAt(chunkX, chunkY)

		g.world.chunks[chunkIdx].grid[tileIdx].Type = 1
	}
}
