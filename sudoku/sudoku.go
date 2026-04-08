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

func IsSolved(puzzle, solution Sudoku) bool {
	for i := range puzzle {
		for j := range puzzle[i] {
			if puzzle[i][j] != solution[i][j] {
				return false
			}
		}
	}
	return true
}

func NewSudoku(level difficulty) (puzzle, solution Sudoku) {
	s := newSudoku()
	fillDiagonal(s)
	fillRemaining(s, 0, 0)
	c := CopyGrid(s)
	removeDigitsUnique(s, level)
	return s, c
}

func ValidAnswer(puzzle, solution Sudoku) (row, col int) {
	for i := range puzzle {
		for j := range puzzle[i] {
			if puzzle[i][j] != 0 && puzzle[i][j] != solution[i][j] {
				return i, j
			}
		}
	}
	return -1, -1
}

func CopyGrid(s Sudoku) Sudoku {
	c := make(Sudoku, size)
	for i := range s {
		c[i] = make([]int, size)
		copy(c[i], s[i])
	}
	return c
}

func newSudoku() Sudoku {
	s := make(Sudoku, size)
	for i := range s {
		s[i] = make([]int, size)
	}
	return s
}

func fillDiagonal(s Sudoku) {
	for i := 0; i < size; i += boxSize {
		fillBox(s, i, i)
	}
}

func fillBox(s Sudoku, row, col int) {
	for i := range boxSize {
		for j := range boxSize {
			for {
				num := rand.Intn(9) + 1
				if unUsedInBox(s, row, col, num) {
					s[row+i][col+j] = num
					break
				}
			}
		}
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

func checkIfSafe(s Sudoku, i, j, num int) bool {
	return unUsedInRow(s, i, num) &&
		unUsedInCol(s, j, num) &&
		unUsedInBox(s, i-i%boxSize, j-j%boxSize, num)
}

func removeDigitsUnique(s Sudoku, k difficulty) {
	for k > 0 {
		i := rand.Intn(size)
		j := rand.Intn(size)

		if s[i][j] == 0 {
			continue
		}

		backup := s[i][j]
		s[i][j] = 0

		tmp := CopyGrid(s)
		if countSolutions(tmp) != 1 {
			s[i][j] = backup
		} else {
			k--
		}
	}
}

func countSolutions(s Sudoku) int {
	count := 0
	var solver func(Sudoku)
	solver = func(board Sudoku) {
		for i := range size {
			for j := range size {
				if board[i][j] == 0 {
					for num := 1; num <= 9; num++ {
						if checkIfSafe(board, i, j, num) {
							board[i][j] = num
							solver(board)
							board[i][j] = 0
						}
					}
					return
				}
			}
		}
		count++
	}
	solver(s)
	return count
}
