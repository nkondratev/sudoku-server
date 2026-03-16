package main

type MoveDTO struct {
	PlayerID int64 `json:"player_id"`
	Row      int   `json:"row"`
	Col      int   `json:"col"`
	Value    int   `json:"value"`
}

type PuzzleDTO struct {
	PlayerID int64   `json:"player_id"`
	Grid     [][]int `json:"grid"`
}
