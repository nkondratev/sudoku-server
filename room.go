package main

import (
	"sudoku-server/sudoku"
	"sync"
)

const countPlayers = 2

type Room struct {
	Mu       sync.Mutex
	players  [countPlayers]*Player
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
