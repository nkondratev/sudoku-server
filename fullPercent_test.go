package main

import (
	"app/sudoku"
	"testing"
)

func TestFullPercent(t *testing.T) {
	p, s := sudoku.NewSudoku(sudoku.Medium)
	t.Log(sudoku.PrettyPrint(p))
	t.Log(sudoku.PrettyPrint(s))
	t.Log(p.FullPercent(s))
}
