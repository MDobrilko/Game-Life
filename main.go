package main

// go build -ldflags -H=windowsgui

import (
	"math/rand"
	"time"

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
	rand.Seed(time.Now().UnixNano())
}

func generateFieldOfCells() [][]bool {
	NCellsX := int64(ScreenWidth / defaultCell.width)
	NCellsY := int64(ScreenHeight / defaultCell.height)

	field := make([][]bool, NCellsX)
	for i := range field {
		field[i] = make([]bool, NCellsY)
	}

	NCells := NCellsX * NCellsY
	// NLivingCells := rand.Int63n(NCells / 2)
	NLivingCells := rand.Int63n(NCells - NCells/4)

	for i := int64(0); i < NLivingCells; i++ {
		newCell := rand.Int63n(NCells)

		x := newCell % NCellsX
		y := newCell / NCellsX

		field[x][y] = true
	}
	return field
}

func isOnField(field [][]bool, x, y int) bool {
	return x >= 0 && y >= 0 && x < len(field) && y < len(field[x])
}

func update(field [][]bool, newField [][]bool) {
	for x, col := range field {
		for y, isAlive := range col {
			counter := 0

			for i := 0; i < len(moveX); i++ {
				if isOnField(field, x+moveX[i], y+moveY[i]) && field[x+moveX[i]][y+moveY[i]] == true {
					counter++
				}
			}

			if counter == 3 || counter == 2 && isAlive {
				newField[x][y] = true
			} else {
				newField[x][y] = false
			}
		}
	}
}

func addCells(newField [][]bool, field [][]bool) {
	for x, column := range field {
		for y, isAlive := range column {
			if isAlive {
				addDefaultCell(pixel.V(
					float64(x)*defaultCell.width,
					float64(y)*defaultCell.height,
				))
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
