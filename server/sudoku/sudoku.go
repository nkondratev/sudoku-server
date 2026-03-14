package sudoku

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

type difficulty int

const (
	Easy   difficulty = 40
	Medium difficulty = 45
	Hard   difficulty = 50
)

func makeGrid(size int) [][]int {
	grid := make([][]int, size)
	for i := range size {
		grid[i] = make([]int, size)
	}
	return grid
}

func CopyGrid(original [][]int) [][]int {
	duplicate := make([][]int, len(original))
	for i := range original {
		duplicate[i] = make([]int, len(original[i]))
		copy(duplicate[i], original[i])
	}
	return duplicate
}

func resolveSudoku(grid [][]int) bool {
	row, col, err := findEmptyCell(grid)
	if err != nil {
		return true
	}

	size := len(grid)
	numbers := rand.Perm(size)

	for _, n := range numbers {
		value := n + 1

		if !validate(grid, row, col, value) {
			continue
		}

		grid[row][col] = value

		if resolveSudoku(grid) {
			return true
		}

		grid[row][col] = 0
	}

	return false
}

func validate(grid [][]int, row, col, value int) bool {
	size := len(grid)

	for i := range size {
		if grid[row][i] == value {
			return false
		}
		if grid[i][col] == value {
			return false
		}
	}

	box := int(math.Sqrt(float64(size)))
	startRow := row - row%box
	startCol := col - col%box

	for i := range box {
		for j := range box {
			if grid[startRow+i][startCol+j] == value {
				return false
			}
		}
	}

	return true
}

func findEmptyCell(grid [][]int) (int, int, error) {
	for r := range grid {
		for c := range grid[r] {
			if grid[r][c] == 0 {
				return r, c, nil
			}
		}
	}

	return 0, 0, errors.New("no empty cell")
}

func PrettyPrint(grid [][]int) {
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == 0 {
				fmt.Printf(". ")
			} else {
				fmt.Printf("%v ", grid[i][j])
			}
		}
		fmt.Printf("\n")
	}
}

func New(size int) [][]int {
	grid := makeGrid(size)
	resolveSudoku(grid)
	return grid
}

func CreatePuzzle(grid [][]int, level difficulty) {
	removeNumbersUnique(grid, int(level))
}

func removeNumbersUnique(grid [][]int, count int) {
	size := len(grid)
	total := size * size

	cells := rand.Perm(total)

	for _, idx := range cells {
		if count <= 0 {
			break
		}

		r := idx / size
		c := idx % size
		backup := grid[r][c]
		grid[r][c] = 0

		if countSolutions(grid, 2) > 1 {
			grid[r][c] = backup
			continue
		}

		count--
	}
}

func countSolutions(grid [][]int, limit int) int {
	row, col, err := findEmptyCell(grid)
	if err != nil {
		return 1
	}

	size := len(grid)
	count := 0

	for v := 1; v <= size; v++ {
		if !validate(grid, row, col, v) {
			continue
		}

		grid[row][col] = v
		count += countSolutions(grid, limit)
		grid[row][col] = 0

		if count >= limit {
			return count
		}
	}

	return count
}
