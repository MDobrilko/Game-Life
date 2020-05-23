package main

import (
	"image/color"
	"math/rand"

	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	_ "golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type cell struct {
	width, height float64
	color         color.Color
}

var defaultCell = cell{10, 10, pixel.RGB(0, 255, 0)}

var imd = imdraw.New(nil)

func addCell(vertex pixel.Vec, width, height float64, color color.Color) {
	imd.Color = colornames.White
	imd.Push(vertex)
	imd.Color = colornames.White
	imd.Push(pixel.V(vertex.X, vertex.Y+height))
	imd.Color = colornames.White
	imd.Push(pixel.V(vertex.X+width, vertex.Y+height))
	imd.Color = colornames.White
	imd.Push(pixel.V(vertex.X+width, vertex.Y))
	imd.Polygon(1)

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

func addDefaultCell(vertex pixel.Vec) {
	addCell(vertex, defaultCell.width, defaultCell.height, defaultCell.color)
}

func addRandColorCell(vertex pixel.Vec) {
	addCell(
		vertex,
		defaultCell.width,
		defaultCell.height,
		colors[rand.Intn(len(colors))],
	)
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
