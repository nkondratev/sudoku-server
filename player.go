package main

import (
	"sudoku-server/sudoku"
	"sync"
	"sync/atomic"
	"time"

	"github.com/olahol/melody"
)

const gameTime = time.Minute * 7

type Player struct {
	Mu          sync.Mutex
	Time        *time.Timer
	Puzzle      sudoku.Sudoku
	Session     *melody.Session
	FillPercent atomic.Int64
}

func NewPlayer(s *melody.Session) *Player {
	return &Player{
		Session: s,
		Time:    time.NewTimer(gameTime),
	}
}
