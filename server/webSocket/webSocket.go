package websocket

import (
	"server/sudoku"
	"sync"
	"sync/atomic"
	"time"

	"github.com/olahol/melody"
)

const (
	gameTime     = 10
	gridSize     = 9
	countPlayers = 2
)

var (
	rooms    []*Room
	playerId atomic.Int64
	roomId   atomic.Int64
	players  = make(map[*melody.Session]*Player)
	mu       sync.Mutex
)

type Player struct {
	id     int64
	lives  int
	name   string
	puzzle sudoku.Sudoku
	conn   *melody.Session
}

func NewPlayer(name string, s *melody.Session, puzzle sudoku.Sudoku) *Player {
	playerId.Add(1)
	copyPuzzle := make(sudoku.Sudoku, len(puzzle))
	for i := range puzzle {
		copyPuzzle[i] = make([]int, len(puzzle[i]))
		copy(copyPuzzle[i], puzzle[i])
	}

	return &Player{
		conn:   s,
		id:     playerId.Load(),
		lives:  3,
		name:   name,
		puzzle: copyPuzzle,
	}
}

type Room struct {
	id       int64
	time     *time.Ticker
	players  []*Player
	solution sudoku.Sudoku
	puzzle   sudoku.Sudoku
}

func NewRoom() *Room {
	s := sudoku.NewSolution(gridSize)
	roomId.Add(1)
	return &Room{
		id:       roomId.Load(),
		players:  make([]*Player, countPlayers),
		solution: s,
		puzzle:   sudoku.NewPuzzle(s, int(sudoku.Easy)),
	}
}
