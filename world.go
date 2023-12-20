package main

import "image/color"

type World struct {
	lastUpdate int64

	grid [][]object
}
type object struct {
	color color.RGBA
}

func (w *World) Update() error {
	return nil
}
