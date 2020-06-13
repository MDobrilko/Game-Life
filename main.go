package main

// go build -ldflags -H=windowsgui

import (
	"fmt"
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

func generateFieldOfDeadCells() [][]bool {
	NCellsX := int64(ScreenWidth / defaultCell.width)
	NCellsY := int64(ScreenHeight / defaultCell.height)

	field := make([][]bool, NCellsX)
	for i := range field {
		field[i] = make([]bool, NCellsY)
	}

	return field
}

func generateFieldOfCells() [][]bool {
	NCellsX := int64(ScreenWidth / defaultCell.width)
	NCellsY := int64(ScreenHeight / defaultCell.height)

	field := generateFieldOfDeadCells()

	NCells := NCellsX * NCellsY
	NLivingCells := rand.Int63n(NCells / 3)
	// NLivingCells := rand.Int63n(NCells - NCells/4)

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

func isGameOver(prevFields [][][]bool, field [][]bool) bool {
	for i := range prevFields {
		areFieldsDifferent := false
	CheckFields:
		for x, col := range prevFields[i] {
			for y, isAlive := range col {
				if isAlive != field[x][y] {
					areFieldsDifferent = true
					break CheckFields
				}
			}
		}
		if areFieldsDifferent == false {
			return true
		}
	}
	return false
}

func parUpdate(field [][]bool, newField [][]bool) {
	goroutinesNum := 200000
	ch := make(chan bool, goroutinesNum)

	for i := 0; i < goroutinesNum; i++ {
		ch <- true
	}

	for x, col := range field {
		<-ch
		go func(field [][]bool, col []bool, x int) {
			for y, isAlive := range col {
				livingCellsNum := 0

				for i := 0; i < len(moveX); i++ {
					if isOnField(field, x+moveX[i], y+moveY[i]) && field[x+moveX[i]][y+moveY[i]] == true {
						livingCellsNum++
					}
				}

				if livingCellsNum == 3 || livingCellsNum == 2 && isAlive {
					newField[x][y] = true
				} else {
					newField[x][y] = false
				}
			}
		}(field, col, x)
		ch <- true
	}
}

func seqUpdate(field [][]bool, newField [][]bool) {
	for x, col := range field {
		for y, isAlive := range col {
			livingCellsNum := 0

			for i := 0; i < len(moveX); i++ {
				if isOnField(field, x+moveX[i], y+moveY[i]) && field[x+moveX[i]][y+moveY[i]] == true {
					livingCellsNum++
				}
			}

			if livingCellsNum == 3 || livingCellsNum == 2 && isAlive {
				newField[x][y] = true
			} else {
				newField[x][y] = false
			}
		}
	}
}

func addCellsToWin(field [][]bool) {
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

func copyField(dst [][]bool, src [][]bool) {
	for i := range dst {
		copy(dst[i], src[i])
	}
}

func addFieldToPrevField(prevFields [][][]bool, field [][]bool) {
	for i := len(prevFields) - 1; i > 0; i-- {
		prevFields[i-1], prevFields[i] = prevFields[i], prevFields[i-1]
	}
	copyField(prevFields[0], field)
}

func startGame(win *pixelgl.Window) {
	prevStepsNum := 8
	var prevFields = make([][][]bool, prevStepsNum)
	for i := 0; i < prevStepsNum; i++ {
		prevFields[i] = generateFieldOfDeadCells()
	}

	field := generateFieldOfCells()
	addCellsToWin(field)
	drawCells(win)

	timer := time.Now()
	frames := 0
	for !isGameOver(prevFields, field) && !win.Closed() {
		win.Clear(colornames.Black)
		clearCells()

		addFieldToPrevField(prevFields, field)
		seqUpdate(prevFields[0], field)
		addCellsToWin(field)

		drawCells(win)
		win.Update()

		frames++
		if time.Now().Sub(timer).Milliseconds() >= 1000 {
			fmt.Println("frames per second = ", frames)
			frames = 0
			timer = time.Now()
		}
	}
	for !win.Closed() {
		win.Update()
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

	startTime := time.Now()
	startGame(win)
	fmt.Println("Time of game: ", time.Now().Sub(startTime).Milliseconds(), "ms")
}

func main() {
	pixelgl.Run(run)
}
