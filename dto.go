package main

import (
	"sudoku-server/sudoku"
	"time"
)

// Это сообщение от клиента к серверу
type MessageDTO struct {
	Puzzle sudoku.Sudoku `json:"puzzle"`
	IsEnd  bool          `json:"is_end"`
}

// Это сообщение от сервера к клиенту
type SendMessageDTO struct {
	IsValid  bool          `json:"is_valid"`
	Row      int           `json:"row"`
	Col      int           `json:"col"`
	Puzzle   sudoku.Sudoku `json:"sudoku"`
	IsSolved bool          `json:"is_solved"`
}

type FirstMessage struct {
	Puzzle sudoku.Sudoku `json:"sudoku"`
	Time   time.Duration `json:"time"`
}
