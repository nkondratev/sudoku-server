package sudoku

import "testing"

func TestSudoku(t *testing.T) {

	s := NewSudoku()
	var grid = *NewSudoku()
	grid.Grid = copyGrid(s.Grid)
	removeNumbersUnique(grid, 10)
}
