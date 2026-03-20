package websocket

import (
	"slices"
	"sync"

	"github.com/olahol/melody"
)

type Server struct {
	Rooms   []*Room
	Players map[*melody.Session]*Player
	mu      sync.Mutex
}

func (s *Server) DeletePlayer(conn *melody.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Players, conn)
}

func NewServer() *Server {
	return &Server{
		Rooms:   make([]*Room, 0, 1024),
		Players: make(map[*melody.Session]*Player),
	}
}

func (s *Server) FindPlayer(conn *melody.Session) (*Player, error) {
	p, ok := s.Players[conn]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

func (s *Server) FindRoom(conn *melody.Session) (*Room, error) {
	p, ok := s.Players[conn]
	if !ok {
		return nil, ErrNotFound
	}
	for i := range s.Rooms {
		if slices.Contains(s.Rooms[i].Players, p) {
			return s.Rooms[i], nil
		}
	}
	return nil, ErrNotFound

}

func (s *Server) Append(room *Room) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Rooms = append(s.Rooms, room)
	for _, player := range room.Players {
		s.Players[player.Conn] = player
	}
}

func (s *Server) Delete(room *Room) error {
	if len(s.Rooms) == 0 {
		return ErrEmpty
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.Rooms {
		if s.Rooms[i] == room {
			s.Rooms = append(s.Rooms[0:i], s.Rooms[i+1:]...)
			return nil
		}
	}

	return ErrNotFound
}
