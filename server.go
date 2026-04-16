package main

import (
	"app/mode"
	"app/sudoku"
	"encoding/json"
	"log/slog"

	"github.com/olahol/melody"
)

type Server struct {
	m        *melody.Melody
	playerCh chan *Player
	logger   *slog.Logger
}

func (s *Server) matchmakingLoop() {
	lobby := make([]*Player, 0, mode.CountPlayers)

	for {
		p := <-s.playerCh
		lobby = append(lobby, p)

		s.logger.Info("player joined lobby", "size", len(lobby))

		if len(lobby) < mode.CountPlayers {
			continue
		}

		// =========================
		// CREATE ROOM
		// =========================

		room := NewRoom()
		s.logger.Info("room created")

		room.players = make([]*Player, mode.CountPlayers)

		for i := 0; i < mode.CountPlayers; i++ {
			room.players[i] = lobby[i]
			room.players[i].Session.Set(RoomString, room)
		}

		// очистить лобби
		lobby = lobby[:0]

		// =========================
		// INIT GAME STATE
		// =========================

		for i := range room.players {
			room.players[i].Puzzle = sudoku.CopyGrid(room.Puzzle)
		}

		init := FirstMessage{
			Puzzle: room.Puzzle,
		}

		data, err := json.Marshal(init)
		if err != nil {
			s.logger.Error("marshal failed", "err", err)
			continue
		}

		for _, p := range room.players {
			_ = p.Session.Write(data)
		}

		s.logger.Info("puzzle sent")
	}
}
