package main

import "sudoku-server/sudoku"

// Это сообщение от клиента к серверу
type MessageDTO struct {
	Puzzle sudoku.Sudoku `json:"puzzle"`
}

// Это сообщение от сервера к клиенту
type SendMessageDTO struct {
	Lives    int           `json:"lives"`
	Row      int           `json:"row"`
	Col      int           `json:"col"`
	Puzzle   sudoku.Sudoku `json:"sudoku"`
	IsSolved bool          `json:"is_solved"`
}
