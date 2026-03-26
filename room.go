package main

import (
	"strconv"
	"sudoku-server/sudoku"
	"sync/atomic"
)

const countPlayers = 2

var id atomic.Int64

type Room struct {
	id       int64
	players  [countPlayers]*Player
	Puzzle   sudoku.Sudoku
	Solution sudoku.Sudoku
}

func NewRoom() *Room {
	p, s := sudoku.NewSudoku(sudoku.Easy)
	return &Room{
		id:       id.Add(1),
		Puzzle:   p,
		Solution: s,
	}
}

func (r *Room) Id() string {
	return strconv.Itoa(int(r.id))
}
