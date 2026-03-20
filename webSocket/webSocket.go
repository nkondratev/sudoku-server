package websocket

import (
	"sudoku-server/sudoku"
	"time"
)

const ()

type Puzzle struct {
	Board    sudoku.Sudoku
	Solution sudoku.Sudoku
}

func HandleRooms(player chan *Player, game chan *Room) {
	for {
		room := NewRoom()
		for i := range room.Players {
			room.Players[i] = <-player
		}
		game <- room
	}
}

func StartGame(r *Room) {
	for _, p := range r.Players {
		if p == nil {
			continue
		}

		p.Puzzle = r.Puzzle.Board

		p.Conn.WebsocketConnection().WriteJSON(map[string]any{
			"puzzle": p.Puzzle,
			"lives":  p.Lives,
		})
		r.StartTimer(time.Minute * sudoku.GameTime)
	}
}
