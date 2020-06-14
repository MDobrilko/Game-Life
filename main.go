package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	ScreenWidth  = 1560.0
	ScreenHeight = 800.0
)

var (
	goroutinesNum = 100
	ch            = make(chan bool, goroutinesNum)

	waitGroup sync.WaitGroup
)

var moveX = [8]int{1, 0, -1, 0, 1, -1, 1, -1}
var moveY = [8]int{0, 1, 0, -1, 1, -1, -1, 1}

func init() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < goroutinesNum; i++ {
		ch <- true
	}
}

func min(val1, val2 int) int {
	if val1 < val2 {
		return val1
	} else {
		return val2
	}
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

func calcNextFieldState(field, newField [][]bool, fromX, fromY, toX, toY int) {
	for x := fromX; x < toX; x++ {
		for y := fromY; y < toY; y++ {
			livingCellsNum := 0

			for i := range moveX {
				if isOnField(field, x+moveX[i], y+moveY[i]) && field[x+moveX[i]][y+moveY[i]] == true {
					livingCellsNum++
				}
			}

			if livingCellsNum == 3 || livingCellsNum == 2 && field[x][y] {
				newField[x][y] = true
			} else {
				newField[x][y] = false
			}
		}
	}
}

func parBlockUpdate(field, newField [][]bool) {
	type Block struct {
		lenY, lenX int
	}

	var (
		fieldLenX   = len(field)
		fieldLenY   = len(field[0])
		blockSquare = (fieldLenX*fieldLenY + goroutinesNum - 1) / goroutinesNum

		blockSize = Block{}
	)

	blockSize.lenX = int(math.Ceil(math.Sqrt(float64(blockSquare) * float64(fieldLenX) / float64(fieldLenY))))
	blockSize.lenY = int(math.Ceil((float64(fieldLenY) / float64(fieldLenX)) * float64(blockSize.lenX)))

	rowBlocksN := (fieldLenX + blockSize.lenX - 1) / blockSize.lenX

	waitGroup.Add(goroutinesNum)
	for idx := 0; idx < goroutinesNum; idx++ {
		go func(field, newField [][]bool, blockSize Block, idx int) {
			fromX := (idx / rowBlocksN) * blockSize.lenX
			toX := min(fromX+blockSize.lenX, fieldLenX)
			fromY := (idx % rowBlocksN) * blockSize.lenY
			toY := min(fromY+blockSize.lenY, fieldLenY)

			if fromX < len(field) && fromY < len(field[0]) {
				calcNextFieldState(field, newField, fromX, fromY, toX, toY)
			}

			waitGroup.Done()
		}(field, newField, blockSize, idx)
	}
	waitGroup.Wait()
}

func parUpdate(field, newField [][]bool) {
	for x, col := range field {
		<-ch
		go func(field [][]bool, col []bool, x int) {
			for y, isAlive := range col {
				livingCellsNum := 0

				for i := range moveX {
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
	calcNextFieldState(field, newField, 0, 0, len(field), len(field[0]))
}

func addCellsToWin(field [][]bool) {
	for x, column := range field {
		for y, isAlive := range column {
			if isAlive {
				addDefaultCell(pixel.V(
					float64(x)*defaultCell.width,
					float64(y)*defaultCell.height,
				))
			} else {
				addDefaultEmptyCell(pixel.V(
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

func showBuilder(win *pixelgl.Window, field [][]bool) {
	newCellX := int(win.MousePosition().X / defaultCell.height)
	newCellY := int(win.MousePosition().Y / defaultCell.width)

	if win.JustPressed(pixelgl.MouseButtonLeft) {
		field[newCellX][newCellY] = true
	} else if win.JustPressed(pixelgl.MouseButtonRight) {
		field[newCellX][newCellY] = false
	}
}

func startGame(cfg *pixelgl.WindowConfig, win *pixelgl.Window) {
	var (
		prevStepsNum = 9
		prevFields   = make([][][]bool, prevStepsNum)
		field        = generateFieldOfDeadCells()

		timer  = time.Tick(time.Second)
		frames = 0

		isFieldGenerated = true
		isLifeGoing      = false
	)

	for i := 0; i < prevStepsNum; i++ {
		prevFields[i] = generateFieldOfDeadCells()
	}

	addCellsToWin(field)
	drawCells(win)

	startTime := time.Now()
	for !win.Closed() {
		win.Clear(colornames.Black)
		clearCells()

		if win.JustPressed(pixelgl.KeyB) || isLifeGoing && isGameOver(prevFields, field) {
			isFieldGenerated = !isFieldGenerated
			isLifeGoing = !isLifeGoing

			fmt.Println("Time of game: ", time.Now().Sub(startTime).Milliseconds(), "ms")
			startTime = time.Now()
		}

		if isLifeGoing {
			addFieldToPrevField(prevFields, field)
			parUpdate(prevFields[0], field)
		} else if isFieldGenerated {
			if win.JustPressed(pixelgl.KeyR) {
				field = generateFieldOfCells()
			}
			showBuilder(win, field)
		}
		addCellsToWin(field)

		drawCells(win)
		win.Update()

		frames++
		select {
		case <-timer:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
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

	startGame(&cfg, win)
}

func main() {
	pixelgl.Run(run)
}
