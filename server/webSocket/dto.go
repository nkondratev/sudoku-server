package websocket

import "sudoku-server/sudoku"

type MoveDTO struct {
	PlayerID int64         `json:"player_id"`
	Sudoku   sudoku.Sudoku `json:"sudoku"`
}
