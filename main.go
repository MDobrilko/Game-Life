package main

// go build -ldflags -H=windowsgui

import (
	"image/color"
	"math/rand"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	ScreenWidth  = 1560.0
	ScreenHeight = 800.0
)

var moveX = [8]int{1, 0, -1, 0, 1, -1, 1, -1}
var moveY = [8]int{0, 1, 0, -1, 1, -1, -1, 1}

func init() {
	rand.Seed(1000)
}

func generateFieldOfCells() [][]*color.RGBA {
	NCellsX := int64(ScreenWidth / defaultCell.width)
	NCellsY := int64(ScreenHeight / defaultCell.height)

	field := make([][]*color.RGBA, NCellsX)
	for i := range field {
		field[i] = make([]*color.RGBA, NCellsY)
	}

	NCells := NCellsX * NCellsY
	// NLivingCells := rand.Int63n(NCells / 2)
	NLivingCells := rand.Int63n(NCells - NCells/4)

	for i := int64(0); i < NLivingCells; i++ {
		newCell := rand.Int63n(NCells)

		x := newCell % NCellsX
		y := newCell / NCellsX

		field[x][y] = &colors[rand.Intn(len(colors))]
	}
	return field
}

func mixColor(otherColors []*color.RGBA) *color.RGBA {
	var mixedColor *color.RGBA = &color.RGBA{0, 0, 0, 255}
	var count uint8
	for _, col := range otherColors {
		if col != nil {
			mixedColor.R += col.R
			mixedColor.G += col.G
			mixedColor.B += col.B
			count++
		}
	}

	if count > 0 {
		mixedColor.B /= count
		mixedColor.G /= count
		mixedColor.R /= count
	}

	var minVal = uint8(40)
	if mixedColor.R < minVal && mixedColor.G < minVal && mixedColor.B < minVal {
		mixedColor.R, mixedColor.G, mixedColor.B, mixedColor.A = minVal, minVal, minVal, 255
	}

	return mixedColor
}

func isOnField(field [][]*color.RGBA, x, y int) bool {
	return x >= 0 && y >= 0 && x < len(field) && y < len(field[x])
}

func update(field [][]*color.RGBA, newField [][]*color.RGBA) {
	var nearbyCells [9]*color.RGBA
	var shadedNearbyCells = nearbyCells[:0]

	for x, col := range field {
		for y, col := range col {
			counter := 0

			for i := 0; i < len(moveX); i++ {
				if isOnField(field, x+moveX[i], y+moveY[i]) && field[x+moveX[i]][y+moveY[i]] != nil {
					shadedNearbyCells = append(shadedNearbyCells, field[x+moveX[i]][y+moveY[i]])
					counter++
				}
			}

			if counter == 3 || counter == 2 && col != nil {
				if len(shadedNearbyCells) > 0 {
					newField[x][y] = mixColor(shadedNearbyCells)
				} else {
					newField[x][y] = &colors[rand.Intn(len(colors))]
				}
			} else {
				newField[x][y] = nil
			}
		}
	}
}

func addCells(newField [][]*color.RGBA, field [][]*color.RGBA) {
	for x, column := range field {
		for y, color := range column {
			if color == nil && newField[x][y] != nil {
				addRandColorCell(pixel.V(
					float64(x)*defaultCell.width,
					float64(y)*defaultCell.height,
				))
			} else if color != nil && newField[x][y] != nil {
				addColorCell(pixel.V(
					float64(x)*defaultCell.width,
					float64(y)*defaultCell.height,
				),
					field[x][y],
				)
			}
		}
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Game life",
		Bounds: pixel.R(0, 0, ScreenWidth, ScreenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	field := generateFieldOfCells()
	newField := generateFieldOfCells()
	addCells(newField, field)
	drawCells(win)

	for !win.Closed() {
		win.Clear(colornames.Black)
		clearCells()

		update(field, newField)
		addCells(newField, field)

		drawCells(win)
		win.Update()

		field, newField = newField, field
		// time.Sleep(150 * time.Millisecond)
	}
}

func main() {
	pixelgl.Run(run)
}
