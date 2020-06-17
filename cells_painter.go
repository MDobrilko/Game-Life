package main

import (
	"image/color"

	"github.com/faiface/pixel/pixelgl"
	_ "golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type cell struct {
	width, height float64
	color         color.Color
}

var defaultCell = cell{1, 1, pixel.RGB(0, 1, 0)}

var imd = imdraw.New(nil)

func addCell(vertex pixel.Vec, width, height float64, color color.Color) {
	imd.Color = color
	imd.Push(vertex)
	imd.Color = color
	imd.Push(pixel.V(vertex.X, vertex.Y+height))
	imd.Color = color
	imd.Push(pixel.V(vertex.X+width, vertex.Y+height))
	imd.Color = color
	imd.Push(pixel.V(vertex.X+width, vertex.Y))
	imd.Polygon(0)
}

func addDefaultEmptyCell(vertex pixel.Vec) {
	imd.Color = pixel.RGB(0.1, 0.1, 0.1)
	imd.Push(vertex)
	imd.Color = pixel.RGB(0.1, 0.1, 0.1)
	imd.Push(pixel.V(vertex.X, vertex.Y+defaultCell.height))
	imd.Color = pixel.RGB(0.1, 0.1, 0.1)
	imd.Push(pixel.V(vertex.X+defaultCell.width, vertex.Y+defaultCell.height))
	imd.Color = pixel.RGB(0.1, 0.1, 0.1)
	imd.Push(pixel.V(vertex.X+defaultCell.width, vertex.Y))
	imd.Polygon(1)
}

func addDefaultCell(vertex pixel.Vec) {
	addCell(vertex, defaultCell.width, defaultCell.height, defaultCell.color)
}

func addColorCell(vertex pixel.Vec, color color.Color) {
	addCell(
		vertex,
		defaultCell.width,
		defaultCell.height,
		color,
	)
}

func drawCells(win *pixelgl.Window) {
	imd.Draw(win)
}

func clearCells() {
	imd.Clear()
}
