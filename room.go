package main

import "sudoku-server/sudoku"

const countPlayers = 2

type Room struct {
	players  [countPlayers]*Player
	Puzzle   sudoku.Sudoku
	Solution sudoku.Sudoku
}

func NewRoom() *Room {
	p, s := sudoku.NewSudoku(sudoku.Easy)
	return &Room{
		Puzzle:   p,
		Solution: s,
		players:  [2]*Player{},
	}
}
