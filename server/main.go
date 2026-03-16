package main

import (
	"encoding/json"
	"log"
	"net/http"
	"server/sudoku"
	"sync"
	"sync/atomic"
	"time"

	"github.com/olahol/melody"
)

const (
	gameTime     = 10
	gridSize     = 9
	countPlayers = 2
)

var (
	rooms    []*Room
	playerId atomic.Int64
	roomId   atomic.Int64
	players  = make(map[*melody.Session]*Player)
	mu       sync.Mutex
)

type Player struct {
	id     int64
	lives  int
	name   string
	puzzle sudoku.Sudoku
	conn   *melody.Session
}

func NewPlayer(name string, s *melody.Session, puzzle sudoku.Sudoku) *Player {
	playerId.Add(1)
	copyPuzzle := make(sudoku.Sudoku, len(puzzle))
	for i := range puzzle {
		copyPuzzle[i] = make([]int, len(puzzle[i]))
		copy(copyPuzzle[i], puzzle[i])
	}

	return &Player{
		conn:   s,
		id:     playerId.Load(),
		lives:  3,
		name:   name,
		puzzle: copyPuzzle,
	}
}

type Room struct {
	id       int64
	time     *time.Ticker
	players  []*Player
	solution sudoku.Sudoku
	puzzle   sudoku.Sudoku
}

func NewRoom() *Room {
	s := sudoku.NewSolution(gridSize)
	roomId.Add(1)
	return &Room{
		id:       roomId.Load(),
		players:  make([]*Player, countPlayers),
		solution: s,
		puzzle:   sudoku.NewPuzzle(s, int(sudoku.Easy)),
	}
}

func handleRoom(ch chan *Player, wg *sync.WaitGroup) {
	defer wg.Done()
	room := NewRoom()
	for i := range room.players {
		room.players[i] = <-ch
	}
	rooms = append(rooms, room)
}

func handleRooms(matchmaking chan *Player) {
	for {
		room := NewRoom()

		for i := range countPlayers {
			room.players[i] = <-matchmaking
		}

		room.time = time.NewTicker(gameTime * time.Minute)
		for _, p := range room.players {
			p.conn.WebsocketConnection().WriteJSON(room.puzzle)
		}
		rooms = append(rooms, room)
		startGame(room)
	}
}

func startGame(room *Room) {

}

func main() {
	var matchmaking = make(chan *Player)
	m := melody.New()

	go handleRooms(matchmaking)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})

	m.HandleConnect(func(s *melody.Session) {
		empty := sudoku.New(gridSize)
		player := NewPlayer("1", s, empty)

		mu.Lock()
		players[s] = player
		mu.Unlock()

		matchmaking <- player
		log.Println("wait player")
	})

	m.HandleDisconnect(func(s *melody.Session) {
		mu.Lock()
		delete(players, s)
		mu.Unlock()
		log.Println("player disconnect")
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var move MoveDTO
		err := json.Unmarshal(msg, &move)
		if err != nil {
			log.Println("invalid move:", err)
			return
		}

		mu.Lock()
		player, ok := players[s]
		mu.Unlock()
		if !ok {
			log.Println("player not found")
			return
		}

		player.puzzle[move.Row][move.Col] = move.Value

		for _, room := range rooms {
			for _, p := range room.players {
				if p == player {
					if move.Value != room.solution[move.Row][move.Col] {
						player.lives--
					}
				}
			}
		}

		log.Printf("Player %d made move: row %d, col %d, value %d, lives %d", player.id, move.Row, move.Col, move.Value, player.lives)
	})
	http.ListenAndServe(":5000", nil)
}
