package main

import (
	"app/sudoku"
	"testing"
)

func TestIsSolved(t *testing.T) {

	p, s := sudoku.NewSudoku(sudoku.Medium)
	value := sudoku.IsSolved(p, s)
	if value {
		t.Error("error")
	}
	t.Log("success")

	value = sudoku.IsSolved(s, s)
	if !value {
		t.Error("error")
	}
	t.Log("success")
}
