package websocket

import (
	"encoding/json"
	"sudoku-server/sudoku"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

var roomId atomic.Int64

type Room struct {
	Id       int64
	mu       *sync.RWMutex
	players  []*Player
	Puzzle   sudoku.Sudoku
	Solution sudoku.Sudoku
}

func NewRoom() *Room {
	p, s := sudoku.NewSudoku(sudoku.Easy)
	return &Room{
		Id:       roomId.Add(1),
		players:  make([]*Player, 0, 2),
		Puzzle:   p,
		Solution: s,
	}
}

func (r *Room) Players() []*Player {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.players
}

func (r *Room) Append(p *Player) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.players = append(r.players, p)
}

func (r *Room) StartGame() {
	if len(r.players) != 2 {
		return
	}
	for _, p := range r.players {
		if p == nil {
			return
		}

		p.Puzzle = r.Puzzle

		data, err := json.Marshal(gin.H{
			"puzzle": p.Puzzle,
			"lives":  p.Lives,
		})
		if err != nil {
			return
		}
		if err := p.Conn.Write(data); err != nil {
			return
		}
	}
}
