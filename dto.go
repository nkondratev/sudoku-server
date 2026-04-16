package main

import (
	"app/sudoku"
)

// Это сообщение от клиента к серверу
type MessageDTO struct {
	Puzzle sudoku.Sudoku `json:"puzzle"`
}

// Это сообщение от сервера к клиенту
type SendMessageDTO struct {
	FullPercent float64 `json:"full_percent"`
	IsSolved    bool    `json:"is_solved"`
}

// Это первое сообщение от сервера к клиенту
type FirstMessage struct {
	Puzzle sudoku.Sudoku `json:"sudoku"`
}
