package main

import (
	"app/sudoku"
	"sync"
	"time"

	"github.com/olahol/melody"
)

const gameTime = time.Minute * 7

type Player struct {
	Mu      sync.Mutex
	Time    *time.Timer
	Puzzle  sudoku.Sudoku
	Session *melody.Session
}

func NewPlayer(s *melody.Session) *Player {
	return &Player{
		Session: s,
		Time:    time.NewTimer(gameTime),
	}
}
