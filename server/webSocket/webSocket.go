package websocket

import (
	"server/sudoku"
	"sync"
	"sync/atomic"
	"time"

	"github.com/olahol/melody"
)

const (
	countPlayers = 2
)

var (
	playerId atomic.Int64
	roomId   atomic.Int64
)

type Player struct {
	id     int64
	lives  int
	puzzle sudoku.Sudoku
	conn   *melody.Session
}

type Puzzle struct {
	board    sudoku.Sudoku
	solution sudoku.Sudoku
}

type Room struct {
	id      int64
	players []*Player
	puzzle  Puzzle
	timer   *time.Timer
}

type Server struct {
	rooms   []*Room
	players map[*melody.Session]*Player
	mu      sync.Mutex
}

func (s *Server) AppendRoom(room *Room) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rooms = append(s.rooms, room)
}

func (s *Server) Rooms() []*Room {
	return s.rooms
}

func (s *Server) Players() map[*melody.Session]*Player {
	return s.players
}

func (s *Server) AppendPlayer(player *Player) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.players[player.conn] = player
}

func (s *Server) Delete(room *Room) error {
	if len(s.rooms) == 0 {
		return ErrEmpty
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.rooms {
		if s.rooms[i] == room {
			s.rooms = append(s.rooms[0:i], s.rooms[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
func (r *Room) StartTimer(t time.Duration) {
	r.timer.Reset(t)
}

func (s *Server) DeletePlayer(player *Player) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.players, player.conn)
}

func NewServer() *Server {
	return &Server{
		rooms:   make([]*Room, 0, 1024),
		players: make(map[*melody.Session]*Player),
	}
}

func (s *Server) FindPlayer(conn *melody.Session) (*Player, error) {
	player, ok := s.players[conn]
	if !ok {
		return nil, ErrNotFound
	}
	return player, nil
}

func HandleRooms(player chan *Player, server *Server) {
	room := NewRoom()
	for i := range room.players {
		room.players[i] = <-player
	}
	server.AppendRoom(room)
	room.StartGame()
}

func NewPlayer(s *melody.Session) *Player {
	return &Player{
		conn:  s,
		id:    playerId.Add(1),
		lives: 3,
	}
}

func (p *Player) SetLives(lives int) {
	p.lives = lives
}

func (p *Player) Lives() int {
	return p.lives
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
		timer: time.NewTimer(sudoku.GameTime * time.Minute),
	}
}
func (r *Room) StartGame() {
	for _, p := range r.players {
		if p == nil {
			continue
		}

		for i := range r.players {
			r.players[i].puzzle = r.puzzle.board
		}

		p.conn.WebsocketConnection().WriteJSON(map[string]any{
			"puzzle": p.puzzle,
			"lives":  p.lives,
		})
	}
}
