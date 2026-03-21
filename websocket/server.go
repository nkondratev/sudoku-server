package websocket

import (
	"sync"

	"github.com/olahol/melody"
)

type Server struct {
	mu    *sync.Mutex
	rooms []*Room
}

func NewServer() *Server {
	return &Server{
		rooms: make([]*Room, 0, 1024),
	}
}

func (s *Server) Append(r *Room) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rooms = append(s.rooms, r)
}

func (s *Server) Delete(r *Room) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.rooms {
		if s.rooms[i].Id == r.Id {
			s.rooms = append(s.rooms[:i], s.rooms[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (s *Server) DisconnectPlayer(session *melody.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.rooms {
		for j := range s.rooms[i].players {
			if s.rooms[i].players[j].Conn == session {
				s.rooms[i].players = append(s.rooms[i].players[:i], s.rooms[i].players[i+1:]...)
				return
			}
		}
	}
}

func (s *Server) FindOrCreateRoom(session *melody.Session) *Room {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := NewPlayer(session)
	for i := range s.rooms {
		if len(s.rooms[i].players) != 2 {
			for j := range s.rooms[i].players {
				if s.rooms[i].players[j] == nil {
					s.rooms[i].players[j] = p
				}
			}
			s.rooms[i].players = append(s.rooms[i].players, p)
			return s.rooms[i]
		}
	}
	r := NewRoom()
	r.Append(p)
	s.Append(r)
	return r
}
