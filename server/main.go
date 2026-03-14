package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"server/sudoku"
	"sync"
)

const (
	size         = 9
	countPlayers = 2
	port         = ":8080"
)

type Player struct {
	Name  string
	Conn  net.Conn
	Lives int
}

type Room struct {
	Players  []*Player
	Puzzle   [][]int
	Solution [][]int
	mu       sync.Mutex
}

func main() {
	players := make(chan *Player, 100)
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	go func() {
		for {
			handleConnection(players)
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go handlePlayer(conn, players)
	}
}
func NewPlayer(conn net.Conn) *Player {
	return &Player{
		Conn:  conn,
		Lives: 3,
	}
}
func handlePlayer(conn net.Conn, ch chan *Player) {
	fmt.Println("new player")
	ch <- NewPlayer(conn)
}

func handleConnection(ch chan *Player) {
	fmt.Println("new room")
	solution := sudoku.New(size)
	puzzle := sudoku.CopyGrid(solution)
	sudoku.CreatePuzzle(puzzle, sudoku.Easy)
	room := &Room{
		Players:  make([]*Player, countPlayers),
		Puzzle:   puzzle,
		Solution: solution,
	}
	sudoku.PrettyPrint(solution)

	for i := range room.Players {
		room.Players[i] = <-ch
	}

	for i := range room.Players {
		json.NewEncoder(room.Players[i].Conn).Encode(room.Puzzle)
	}
	handleGame(room)

}
func handleGame(room *Room) {
	var wg sync.WaitGroup

	for _, player := range room.Players {
		wg.Add(1)
		go func(p *Player) {
			defer wg.Done()
			for {
				var board [][]int
				err := json.NewDecoder(p.Conn).Decode(&board)
				if err != nil {
					log.Println("Player disconnected:", err)
					p.Conn.Close()
					return
				}

				errors := checkErrors(board, room.Solution)
				if errors > 0 {
					p.Lives -= errors
					fmt.Println("Player lives:", p.Lives)
					if p.Lives <= 0 {
						fmt.Println("Player out of lives:", p.Name)
						// TODO
						return
					}
				}

				resp := map[string]any{
					"board":  board,
					"lives":  p.Lives,
					"errors": errors,
				}
				json.NewEncoder(p.Conn).Encode(resp)
			}
		}(player)
	}

	wg.Wait()
}

func checkErrors(board, solution [][]int) int {
	errors := 0
	for i := range size {
		for j := range size {
			if board[i][j] != 0 && board[i][j] != solution[i][j] {
				errors++
			}
		}
	}
	return errors
}
