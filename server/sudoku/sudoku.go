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

	size         = 9
	GameTime     = 10
	CountPlayers = 2
)

type Sudoku struct {
	Grid [][]int
	Size int
}

func NewSudoku() *Sudoku {
	return &Sudoku{
		Grid: makeGrid(),
		Size: 9,
	}
}

func (s Sudoku) NewPuzzle(level difficulty) Sudoku {
	var grid = *NewSudoku()
	grid.Grid = copyGrid(s.Grid)
	removeNumbersUnique(grid, level)
	return grid
}

func (s Sudoku) NewSolution() Sudoku {
	resolveSudoku(s)
	return s
}

func (s *Sudoku) Compare(grid Sudoku) (mistakes int) {
	for i := range s.Grid {
		for j := range s.Grid[i] {
			if s.Grid[i][j] != 0 && s.Grid[i][j] != grid.Grid[i][j] {
				mistakes++
			}
		}
	}
	return
}

func makeGrid() [][]int {
	grid := make([][]int, 9)
	for i := range grid {
		grid[i] = make([]int, 9)
	}
	return grid
}

func copyGrid(original [][]int) [][]int {
	duplicate := make([][]int, len(original))
	for i := range original {
		duplicate[i] = make([]int, len(original[i]))
		copy(duplicate[i], original[i])
	}
	return duplicate
}

func resolveSudoku(grid Sudoku) bool {
	row, col, err := findEmptyCell(grid)
	if err != nil {
		return true
	}

	numbers := rand.Perm(size)

	for _, n := range numbers {
		value := n + 1

		if !validate(grid.Grid, row, col, value) {
			continue
		}

		grid.Grid[row][col] = value

		if resolveSudoku(grid) {
			return true
		}

		grid.Grid[row][col] = 0
	}

	return false
}

func validate(grid [][]int, row, col, value int) bool {
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

func findEmptyCell(grid Sudoku) (int, int, error) {
	for r := range grid.Grid {
		for c := range grid.Grid[r] {
			if grid.Grid[r][c] == 0 {
				return r, c, nil
			}
		}
	}

	return 0, 0, errors.New("no empty cell")
}

func PrettyPrint(grid Sudoku) {
	box := 3

	for i := range size {
		if i%box == 0 && i != 0 {
			fmt.Println("------+-------+------")
		}

		for j := range size {
			if j%box == 0 && j != 0 {
				fmt.Print("| ")
			}

			if grid.Grid[i][j] == 0 {
				fmt.Print(". ")
			} else {
				fmt.Printf("%d ", grid.Grid[i][j])
			}
		}
		fmt.Println()
	}
}

func removeNumbersUnique(grid Sudoku, count difficulty) {
	total := size * size

	cells := rand.Perm(total)

	for _, idx := range cells {
		if count <= 0 {
			break
		}

		r := idx / size
		c := idx % size
		backup := grid.Grid[r][c]
		grid.Grid[r][c] = 0

		if countSolutions(grid, 2) > 1 {
			grid.Grid[r][c] = backup
			continue
		}

		count--
	}
}

func countSolutions(grid Sudoku, limit int) int {
	row, col, err := findEmptyCell(grid)
	if err != nil {
		return 1
	}

	count := 0

	for v := 1; v <= size; v++ {
		if !validate(grid.Grid, row, col, v) {
			continue
		}

		grid.Grid[row][col] = v
		count += countSolutions(grid, limit)
		grid.Grid[row][col] = 0

		if count >= limit {
			return count
		}
	}

	return count
}
