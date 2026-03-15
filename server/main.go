package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"server/sudoku"
	"sync"
	"time"
)

const (
	size         = 9
	countPlayers = 2
	gameTime     = 10
)

type Player struct {
	Name        string
	Conn        net.Conn
	Lives       int
	playerBoard sync.Map
}

type Room struct {
	Players  []*Player
	Puzzle   [][]int
	Solution [][]int
	timer    *time.Ticker
}

func main() {
	secure := flag.Bool("secure", true, "use for enable or disable presetting ip")
	flag.Parse()
	var ip = []byte(":8080")
	if !*secure {
		resp, err := http.Get("https://ipinfo.io/ip")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		ip, err = io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		panic("dont secure")
	}

	fmt.Println(string(ip))
	players := make(chan *Player, 100)

	ln, err := net.Listen("tcp", string(ip))
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	go roomManager(players)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}

		go handlePlayer(conn, players)
	}
}

func roomManager(players chan *Player) {
	for {
		handleConnection(players)
	}
}

func NewPlayer(conn net.Conn) *Player {
	return &Player{
		Conn:        conn,
		Lives:       3,
		playerBoard: sync.Map{},
	}
}

func handlePlayer(conn net.Conn, ch chan *Player) {
	fmt.Println("New player connected")
	ch <- NewPlayer(conn)
}

func handleConnection(ch chan *Player) {
	fmt.Println("Creating new room")

	solution := sudoku.New(size)
	puzzle := sudoku.CopyGrid(solution)
	sudoku.CreatePuzzle(puzzle, sudoku.Easy)

	room := &Room{
		Players:  make([]*Player, countPlayers),
		Puzzle:   puzzle,
		Solution: solution,
	}
	sudoku.PrettyPrint(room.Solution)

	for i := range room.Players {
		room.Players[i] = <-ch
	}

	for _, p := range room.Players {
		json.NewEncoder(p.Conn).Encode(room.Puzzle)
	}

	room.timer = time.NewTicker(gameTime * time.Minute)

	handleGame(room)
}

func savePlayerBoard(p *Player, board [][]int) {
	copied := sudoku.CopyGrid(board)
	p.playerBoard.Store(p, copied)
}

func getPlayerBoard(p *Player) [][]int {
	if v, ok := p.playerBoard.Load(p); ok {
		return v.([][]int)
	}
	return make([][]int, size)
}

func countCorrectCells(board, solution [][]int) int {
	count := 0
	for i := range board {
		for j := range board[i] {
			if board[i][j] == solution[i][j] {
				count++
			}
		}
	}
	return count
}

func handleGame(room *Room) {
	var wg sync.WaitGroup

	for _, player := range room.Players {
		wg.Add(1)

		go func(p *Player) {
			defer wg.Done()

			for {
				select {
				case <-room.timer.C:
					fmt.Println("Time is up")

					results := make(map[string]int)
					for _, pl := range room.Players {
						results[pl.Name] = countCorrectCells(getPlayerBoard(pl), room.Solution)
					}

					var winner string
					max := -1
					for name, correct := range results {
						if correct > max {
							max = correct
							winner = name
						}
					}

					resp := map[string]any{
						"winner": winner,
						"scores": results,
					}

					for _, pl := range room.Players {
						json.NewEncoder(pl.Conn).Encode(resp)
						pl.Conn.Close()
					}

					return

				default:
				}

				var board [][]int
				p.Conn.SetReadDeadline(time.Now().Add(time.Second))
				err := json.NewDecoder(p.Conn).Decode(&board)
				if err != nil {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						continue
					}
					log.Println("Player disconnected:", err)
					p.Conn.Close()
					return
				}

				errors, row, column := checkErrors(board, room.Solution)
				if errors > 0 {
					p.Lives -= errors
					fmt.Println("Player lives:", p.Lives)
					if p.Lives <= 0 {
						fmt.Println("Player lost:", p.Name)
						return
					}
				}

				savePlayerBoard(p, board)

				resp := map[string]any{
					"row":    row,
					"column": column,
					"lives":  p.Lives,
					"errors": errors,
				}
				json.NewEncoder(p.Conn).Encode(resp)
			}
		}(player)
	}

	wg.Wait()
}

func checkErrors(board, solution [][]int) (int, int, int) {
	var row, column int
	errors := 0

	for i := range size {
		for j := range size {
			if board[i][j] != 0 && board[i][j] != solution[i][j] {
				errors++
				row = i
				column = j
			}
		}
	}

	return errors, row, column
}
