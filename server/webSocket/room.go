package websocket

import (
	"sudoku-server/sudoku"
	"sync/atomic"
	"time"
)

var roomId atomic.Int64

type Room struct {
	id      int64
	players []*Player
	puzzle  Puzzle
	timer   *time.Timer
}

func NewRoom() *Room {
	s := sudoku.NewSudoku()
	return &Room{
		id:      roomId.Add(1),
		players: make([]*Player, countPlayers),
		puzzle: Puzzle{
			board:    s.NewPuzzle(sudoku.Easy),
			solution: s.NewSolution(),
		},
	}
}
func (r *Room) StartTimer(t time.Duration) {
	r.timer.Reset(t)
}

func (r *Room) Players() []*Player {
	return r.players
}
