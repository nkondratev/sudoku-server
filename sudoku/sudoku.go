package sudoku

import (
	"math/rand"
	"slices"
)

type difficulty int

const (
	Easy   difficulty = 40
	Medium difficulty = 45
	Hard   difficulty = 50

	size     = 9
	boxSize  = 3
	GameTime = 10
)

type Sudoku [][]int

func unUsedInBox(s Sudoku, row, col, num int) bool {
	for i := range boxSize {
		for j := range boxSize {
			if s[row+i][col+j] == num {
				return false
			}
		}
	}
	return true
}

func fillBox(s Sudoku, row, col int) {
	var num int
	for i := range boxSize {
		for j := range boxSize {
			for {
				num = rand.Intn(9) + 1
				if unUsedInBox(s, i, j, num) {
					break
				}
			}
			s[row+i][col+j] = num
		}
	}
}

func unUsedInRow(s Sudoku, i, num int) bool {
	return !slices.Contains(s[i], num)
}

func unUsedInCol(s Sudoku, j, num int) bool {
	for i := range size {
		if s[i][j] == num {
			return false
		}
	}
	return true
}

func checkIfSafe(s Sudoku, i, j, num int) bool {
	return unUsedInRow(s, i, num) &&
		unUsedInCol(s, j, num) &&
		unUsedInBox(s, i-i%boxSize, j-j%boxSize, num)
}

func fillDiagonal(s Sudoku) {
	for i := 0; i < size; i += 3 {
		fillBox(s, i, i)
	}
}

func fillRemaining(s Sudoku, i, j int) bool {
	if i == size {
		return true
	}

	if j == size {
		return fillRemaining(s, i+1, 0)
	}

	if s[i][j] != 0 {
		return fillRemaining(s, i, j+1)
	}

	for num := 1; num <= size; num++ {
		if checkIfSafe(s, i, j, num) {
			s[i][j] = num
			if fillRemaining(s, i, j+1) {
				return true
			}
			s[i][j] = 0
		}
	}
	return false
}

func removeDigits(s Sudoku, k difficulty) {
	for k > 0 {
		cellId := rand.Intn(81)

		i := cellId / 9

		j := cellId % 9

		if s[i][j] != 0 {
			s[i][j] = 0
			k -= 1
		}
	}
}

func NewSudoku(level difficulty) (puzzle, solution Sudoku) {
	s := newSudoku()
	fillDiagonal(s)
	fillRemaining(s, 0, 0)
	c := copyGrid(s)
	removeDigits(s, level)
	return s, c

}

func copyGrid(s Sudoku) Sudoku {
	c := make(Sudoku, size)
	for i := range s {
		c[i] = make([]int, size)
		copy(c[i], s[i])
	}
	return c
}

func newSudoku() Sudoku {
	s := make(Sudoku, size)
	for i := range size {
		s[i] = make([]int, size)
	}
	return s
}
