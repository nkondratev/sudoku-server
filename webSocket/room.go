package websocket

import (
	"sudoku-server/sudoku"
	"sync/atomic"
	"time"
)

var roomId atomic.Int64

type Room struct {
	Id      int64
	Players []*Player
	Puzzle  Puzzle
	Timer   *time.Timer
}

func NewRoom() *Room {
	s := sudoku.NewSudoku()
	return &Room{
		Id:      roomId.Add(1),
		Players: make([]*Player, sudoku.CountPlayers),
		Puzzle: Puzzle{
			Board:    s.NewPuzzle(sudoku.Easy),
			Solution: s.NewSolution(),
		},
	}
}
func (r *Room) StartTimer(t time.Duration) {
	r.Timer.Reset(t)
}
