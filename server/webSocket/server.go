package websocket

import (
	"slices"
	"sync"

	"github.com/olahol/melody"
)

type Server struct {
	rooms   []*Room
	players map[*melody.Session]*Player
	mu      sync.Mutex
}

func (s *Server) DeletePlayer(conn *melody.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.players, conn)
}

func NewServer() *Server {
	return &Server{
		rooms:   make([]*Room, 0, 1024),
		players: make(map[*melody.Session]*Player),
	}
}

func (s *Server) FindPlayer(conn *melody.Session) (*Player, error) {
	p, ok := s.players[conn]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

func (s *Server) FindRoom(conn *melody.Session) (*Room, error) {
	p, ok := s.players[conn]
	if !ok {
		return nil, ErrNotFound
	}
	for i := range s.rooms {
		if slices.Contains(s.rooms[i].players, p) {
			return s.rooms[i], nil
		}
	}
	return nil, ErrNotFound

}

func (s *Server) Append(room *Room) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rooms = append(s.rooms, room)
	for _, player := range room.players {
		s.players[player.Conn] = player
	}
}

func (s *Server) Rooms() []*Room {
	return s.rooms
}

func (s *Server) Players() map[*melody.Session]*Player {
	return s.players
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
