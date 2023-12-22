package main

import "image/color"

type World struct {
	lastUpdate int64

	chunks []loadedChunk
}
type object struct {
	color color.RGBA
}

type loadedChunk struct {
	x int
	y int

	grid [256]object
}

func (w *World) Update() error {
	return nil
}
