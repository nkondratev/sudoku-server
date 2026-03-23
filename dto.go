package main

import (
	"sudoku-server/sudoku"
)

type SendMesageDTO struct {
	Puzzle sudoku.Sudoku `json:"puzzle"`
}

type GetMessageDTO struct {
	Lives  int           `json:"lives"`
	Row    int           `json:"row"`
	Col    int           `json:"col"`
	Puzzle sudoku.Sudoku `json:"sudoku"`
}
