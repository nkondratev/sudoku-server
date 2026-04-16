package main

import (
	"app/sudoku"
	"log/slog"
	"sync"
)

type Room struct {
	Mu       sync.Mutex
	players  []*Player
	Puzzle   [][]int
	Solution [][]int
	Closed   bool
}

func (r *Room) Close(logger *slog.Logger) {
	r.Mu.Lock()
	if r.Closed {
		r.Mu.Unlock()
		return
	}
	r.Closed = true
	players := r.players
	r.Mu.Unlock()

	for _, p := range players {
		_ = p.Session.Close()
	}

	logger.Info("room closed")
}
func NewRoom() *Room {
	p, s := sudoku.NewSudoku(sudoku.Medium)
	return &Room{
		Puzzle:   p,
		Solution: s,
		players:  make([]*Player, 2),
	}
}
