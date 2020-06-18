package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	// ScreenWidth  = 1560.0
	// ScreenHeight = 800.0

	ScreenWidth  = 500.0
	ScreenHeight = 500.0
)

var (
	goroutinesNum = 1000
	ch            = make(chan bool, goroutinesNum)
	waitGroup     sync.WaitGroup
)

var moveX = [8]int{1, 0, -1, 0, 1, -1, 1, -1}
var moveY = [8]int{0, 1, 0, -1, 1, -1, -1, 1}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func min(val1, val2 int) int {
	if val1 < val2 {
		return val1
	} else {
		return val2
	}
}

func readField(filename string) [][]bool {
	field := generateFieldOfDeadCells()

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("File hasnt been opened")
		return field
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for i := 0; i < len(field); i++ {
		if !scanner.Scan() {
			fmt.Println("Field in file is to small")
			break
		}
		text := scanner.Text()
		for j := 0; j < len(field[i]); j++ {
			if j > len(text) {
				break
			}
			if text[j] == '1' {
				field[i][j] = true
			}
		}
	}

	return field
}

func writeField(filename string, field [][]bool) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("File hasnt been opened")
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, col := range field {
		for _, isAlive := range col {
			if isAlive {
				writer.WriteString("1")
			} else {
				writer.WriteString("0")
			}
		}
		writer.WriteString("\n")
	}
}

func readFilenameFromConsole() string {
	var filename string
	fmt.Scan(&filename)
	return filename + ".txt"
}

func clearField(field [][]bool) {
	for x, col := range field {
		for y := range col {
			field[x][y] = false
		}
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

func printHelpToConsole() {
	fmt.Println(
		`p - pause and continue game
s - save current field to file "field.txt"
f - read field from file "field.txt"
o - read field from custom file.txt
c - clear field
h - showing current helps
`,
	)
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

func copyField(dst, src [][]bool) {
	for i := range dst {
		copy(dst[i], src[i])
	}
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
		prevStepsNum = 8
		prevFields   = make([][][]bool, prevStepsNum)
		field        = generateFieldOfDeadCells()

		timer  = time.Tick(time.Second)
		frames = 0

		isFieldGenerated = true
		isLifeGoing      = false

		filename = "field.txt"
	)
	fmt.Printf("Field size: %d %d\n\n", len(field), len(field[0]))

	for i := 0; i < prevStepsNum; i++ {
		prevFields[i] = generateFieldOfDeadCells()
	}

	addCellsToWin(field)
	drawCells(win)

	startTime := time.Now()
	for !win.Closed() {
		win.Clear(colornames.Black)
		clearCells()

		if win.JustPressed(pixelgl.KeyP) || isLifeGoing && isGameOver(prevFields, field) {
			isFieldGenerated = !isFieldGenerated
			isLifeGoing = !isLifeGoing

			fmt.Println("Time of game: ", time.Now().Sub(startTime).Milliseconds(), "ms")
			startTime = time.Now()
		}

		if isLifeGoing {
			addFieldToPrevField(prevFields, field)
			parStripUpdate(prevFields[0], field)
		} else if isFieldGenerated {
			if win.JustPressed(pixelgl.KeyR) {
				field = generateFieldOfCells()
			} else if win.JustPressed(pixelgl.KeyS) {
				writeField(filename, field)
			} else if win.JustPressed(pixelgl.KeyF) {
				field = readField(filename)
			} else if win.JustPressed(pixelgl.KeyC) {
				clearField(field)
			} else if win.JustPressed(pixelgl.KeyO) {
				fmt.Print("Insert filename with field: ")
				field = readField(readFilenameFromConsole())
			} else if win.JustPressed(pixelgl.KeyH) {
				printHelpToConsole()
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

func checkField(filename string, field [][]bool) {
	checkingField := readField(filename)
	for i, col := range field {
		for j := range col {
			if field[i][j] != checkingField[i][j] {
				fmt.Printf("Cells are different in {%d, %d}\n", i, j)
			}
		}
	}
}

func addFieldToPrevField(prevFields [][][]bool, field [][]bool) {
	for i := len(prevFields) - 1; i > 0; i-- {
		prevFields[i-1], prevFields[i] = prevFields[i], prevFields[i-1]
	}
	// copyField(prevFields[0], field)
	parCopyField(prevFields[0], field)
}

func startGameWithoutWin() {
	var (
		prevStepsNum = 9
		prevFields   = make([][][]bool, prevStepsNum)

		filename       = "field.txt"
		outputFilename = "output.txt"

		field = readField(filename)
		// field = generateFieldOfCells()
	)
	// writeField(filename, field)
	fmt.Printf("Field size: %d %d\n\n", len(field), len(field[0]))

	for i := 0; i < prevStepsNum; i++ {
		prevFields[i] = generateFieldOfDeadCells()
	}

	startTime := time.Now()
	for !parIsGameOver(prevFields, field) {
		addFieldToPrevField(prevFields, field)
		// seqUpdate(prevFields[0], field)
		// parUpdate(prevFields[0], field)
		parStripUpdate(prevFields[0], field)
	}
	fmt.Println("Time of game: ", time.Now().Sub(startTime).Milliseconds(), "ms")
	// checkField(outputFilename, field)
}

func run() {
	/*cfg := pixelgl.WindowConfig{
		Title:  "Game life",
		Bounds: pixel.R(0, 0, ScreenWidth, ScreenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}*/

	// startGame(&cfg, win)
	startGameWithoutWin()
}

func main() {
	pixelgl.Run(run)
	close(ch)
}
