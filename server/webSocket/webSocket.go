package websocket

import (
	"sudoku-server/sudoku"
	"time"
)

const (
	countPlayers = 2
)

type Puzzle struct {
	board    sudoku.Sudoku
	solution sudoku.Sudoku
}

func HandleRooms(player chan *Player, game chan *Room) {
	for {
		room := NewRoom()
		for i := range room.players {
			room.players[i] = <-player
		}
		game <- room
	}
}

func StartGame(r *Room) {
	for _, p := range r.players {
		if p == nil {
			continue
		}

		p.Puzzle = r.puzzle.board

		p.Conn.WebsocketConnection().WriteJSON(map[string]any{
			"puzzle": p.Puzzle,
			"lives":  p.Lives,
		})
		r.StartTimer(time.Minute * sudoku.GameTime)
	}
}
