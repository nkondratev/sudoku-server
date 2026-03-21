package websocket

import (
	"sudoku-server/sudoku"
	"sync/atomic"

	"github.com/olahol/melody"
)

var playerId atomic.Int64

type Player struct {
	Id     int64
	Lives  int
	Puzzle sudoku.Sudoku
	Conn   *melody.Session
}

func NewPlayer(s *melody.Session) *Player {
	return &Player{
		Conn:  s,
		Id:    playerId.Add(1),
		Lives: 3,
	}
}
