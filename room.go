package main

import (
	"app/mode"
	"app/sudoku"
	"sync"
)

type Room struct {
	Mu       sync.Mutex
	players  [mode.CountPlayers]*Player
	Puzzle   sudoku.Sudoku
	Solution sudoku.Sudoku
	Closed   bool
}

func NewRoom() *Room {
	p, s := sudoku.NewSudoku(sudoku.Easy)
	return &Room{
		Puzzle:   p,
		Solution: s,
	}
}
