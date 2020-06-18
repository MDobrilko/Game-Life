package main

import "sync"

func parCalcNextFieldState(field, newField [][]bool, fromX, fromY, toX, toY int, ch chan bool, waitGroup *sync.WaitGroup) {
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
	<-ch
	waitGroup.Done()
}

func parCopyField(dst, src [][]bool) {
	for i := range dst {
		ch <- true
		waitGroup.Add(1)
		go func(dst, src [][]bool, i int, ch chan bool, waitGroup *sync.WaitGroup) {
			copy(dst[i], src[i])
			<-ch
			waitGroup.Done()
		}(dst, src, i, ch, &waitGroup)
	}
	waitGroup.Wait()
}

func parStripUpdate(field, newField [][]bool) {
	stripSize := len(field) / goroutinesNum

	for i := 0; i < len(field); i += stripSize {
		ch <- true
		waitGroup.Add(1)
		go parCalcNextFieldState(field, newField, i, 0, i+stripSize, len(field[i]), ch, &waitGroup)
	}
	waitGroup.Wait()
}

func parUpdate(field, newField [][]bool) {
	for x := range field {
		ch <- true
		waitGroup.Add(1)
		go parCalcNextFieldState(field, newField, x, 0, x+1, len(field[x]), ch, &waitGroup)
		// go func(field [][]bool, col []bool, x int, ch chan bool, waitGroup *sync.WaitGroup) {
		// 	for y, isAlive := range col {
		// 		livingCellsNum := 0

		// 		for i := range moveX {
		// 			if isOnField(field, x+moveX[i], y+moveY[i]) && field[x+moveX[i]][y+moveY[i]] == true {
		// 				livingCellsNum++
		// 			}
		// 		}

		// 		if livingCellsNum == 3 || livingCellsNum == 2 && isAlive {
		// 			newField[x][y] = true
		// 		} else {
		// 			newField[x][y] = false
		// 		}
		// 	}
		// 	<-ch
		// 	waitGroup.Done()
		// }(field, col, x, ch, &waitGroup)
	}
	waitGroup.Wait()
}

func parIsGameOver(prevFields [][][]bool, field [][]bool) bool {
	isOver := false

	for i := range prevFields {
		areFieldsDifferent := false

		for x, col := range prevFields[i] {
			ch <- true
			waitGroup.Add(1)

			go func(field [][]bool, col []bool, x int, areFieldsDifferent *bool, ch chan bool, waitGroup *sync.WaitGroup, m *sync.Mutex) {
				local := false
				for y, isAlive := range col {
					if isAlive != field[x][y] {
						local = true
						break
					}
				}
				m.Lock()
				*areFieldsDifferent = local
				m.Unlock()
				<-ch
				waitGroup.Done()
			}(field, col, x, &areFieldsDifferent, ch, &waitGroup, m)

			if areFieldsDifferent {
				break
			}
		}
		if areFieldsDifferent == false {
			isOver = true
			break
		}
	}

	waitGroup.Wait()
	return isOver
}
