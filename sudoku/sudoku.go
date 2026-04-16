package sudoku

import (
	"math/rand"
)

type difficulty int

const (
	Easy   difficulty = 30
	Medium difficulty = 40
	Hard   difficulty = 50

	size       = 9
	boxSize    = 3
	countCells = 81
)

type Sudoku [][]int

// =========================
// PUBLIC API
// =========================

func NewSudoku(level difficulty) (puzzle, solution Sudoku) {
	grid := newSudoku()

	fillGrid(grid)

	solution = CopyGrid(grid)

	removeNumbers(grid, int(level))

	return grid, solution
}

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

func ValidAnswer(puzzle, solution Sudoku) bool {
	for i := range puzzle {
		for j := range puzzle[i] {
			if puzzle[i][j] != 0 && puzzle[i][j] != solution[i][j] {
				return false
			}
		}
	}
	return true
}

// =========================
// CORE LOGIC (как в C#)
// =========================

func fillGrid(grid Sudoku) bool {
	for row := range size {
		for col := range size {
			if grid[row][col] == 0 {

				numbers := shuffledNumbers()

				for _, num := range numbers {
					if isValid(grid, row, col, num) {
						grid[row][col] = num

						if fillGrid(grid) {
							return true
						}

						grid[row][col] = 0
					}
				}

				return false
			}
		}
	}
	return true
}

func shuffledNumbers() []int {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	for i := range nums {
		j := rand.Intn(len(nums)-i) + i
		nums[i], nums[j] = nums[j], nums[i]
	}

	return nums
}

// =========================
// REMOVE NUMBERS (уникальность)
// =========================

func removeNumbers(grid Sudoku, count int) {
	for count > 0 {
		row := rand.Intn(size)
		col := rand.Intn(size)

		if grid[row][col] == 0 {
			continue
		}

		backup := grid[row][col]
		grid[row][col] = 0

		tmp := CopyGrid(grid)
		if countSolutions(tmp) != 1 {
			grid[row][col] = backup
		} else {
			count--
		}
	}
}

// =========================
// COUNT SOLUTIONS
// =========================

func countSolutions(grid Sudoku) int {
	count := 0

	var solve func(Sudoku)
	solve = func(board Sudoku) {
		if count > 1 {
			return
		}

		for row := range size {
			for col := range size {
				if board[row][col] == 0 {
					for num := 1; num <= 9; num++ {
						if isValid(board, row, col, num) {
							board[row][col] = num
							solve(board)
							board[row][col] = 0
						}
					}
					return
				}
			}
		}

		count++
	}

	solve(grid)
	return count
}

// =========================
// VALIDATION
// =========================

func isValid(grid Sudoku, row, col, num int) bool {
	for i := range size {
		if grid[row][i] == num || grid[i][col] == num {
			return false
		}
	}

	startRow := (row / boxSize) * boxSize
	startCol := (col / boxSize) * boxSize

	for r := startRow; r < startRow+boxSize; r++ {
		for c := startCol; c < startCol+boxSize; c++ {
			if grid[r][c] == num {
				return false
			}
		}
	}

	return true
}

// =========================
// HELPERS
// =========================

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

func (s Sudoku) FullPercent(solution Sudoku) float64 {
	var countValid float64
	for i := range s {
		for j := range s[i] {
			if s[i][j] == solution[i][j] {
				countValid++
			}
		}
	}
	percent := countValid / countCells * 100
	return percent
}
