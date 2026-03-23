package main

import (
	"sudoku-server/sudoku"

	"github.com/olahol/melody"
)

const countLives = 3

type Player struct {
	Lives   int
	Puzzle  sudoku.Sudoku
	Session *melody.Session
}

func NewPlayer(s *melody.Session) *Player {
	return &Player{
		Session: s,
		Lives:   countLives,
	}
}
