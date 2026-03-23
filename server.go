package main

import (
	"sync"
)

type Server struct {
	mu    *sync.Mutex
	rooms []*Room
	join  chan *Room
}

func NewServer() *Server {
	return &Server{
		rooms: make([]*Room, 1024),
	}
}
