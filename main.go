package main

import (
	"app/sudoku"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

const (
	addr         = ":8080"
	RoomString   = "room"
	PlayerString = "player"
)

func main() {

	logDir := "./logs/"
	logFile := filepath.Join(logDir, "logs")

	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		log.Fatalf("failed to create log dir: %v", err)
	}

	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer f.Close()
	stream := io.MultiWriter(os.Stdout, f)

	logger := log.New(stream, "", log.Ldate|log.Ltime)

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	m := melody.New()
	playerCh := make(chan *Player)

	go func() {
		var room *Room
		for {

			room = NewRoom()

			for i := range room.players {
				room.players[i] = <-playerCh
				room.players[i].Session.Set(RoomString, room)
			}

			for i := range room.players {
				room.players[i].Puzzle = sudoku.CopyGrid(room.Puzzle)
			}

			fm := &FirstMessage{
				Puzzle: room.Puzzle,
			}

			jsonData, _ := json.Marshal(fm)
			for i := range room.players {
				room.players[i].Session.Write(jsonData)
			}
			logger.Println("data is sented to player")
		}
	}()

	m.HandleConnect(func(s *melody.Session) {
		player := NewPlayer(s)
		logger.Println("new player connected")

		s.Set(PlayerString, player)
		go func() { playerCh <- player }()
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {

		player := s.MustGet(PlayerString).(*Player)
		room := s.MustGet(RoomString).(*Room)

		msgDTO := &MessageDTO{}
		json.Unmarshal(msg, msgDTO)

		if msgDTO.IsEnd {
			room.Mu.Lock()
			if room.Closed {
				room.Mu.Unlock()
				logger.Println("room is closed")
				return
			}
			room.Closed = true
			players := room.players
			room.Mu.Unlock()

			for _, p := range players {
				p.Session.Close()
			}
			return
		}

		player.Mu.Lock()
		player.Puzzle = msgDTO.Puzzle
		puzzleCopy := sudoku.CopyGrid(player.Puzzle)
		player.Mu.Unlock()

		solved := sudoku.IsSolved(puzzleCopy, room.Solution)
		valid := sudoku.ValidAnswer(puzzleCopy, room.Solution)

		var secondPlayer *Player
		for _, p := range room.players {
			if p != player {
				secondPlayer = p
				break
			}
		}

		if secondPlayer == nil {
			return
		}

		// Получаем пазл второго игроков
		secondPlayer.Mu.Lock()
		secondPuzzle := sudoku.CopyGrid(secondPlayer.Puzzle)
		secondPlayer.Mu.Unlock()

		sendmsgDTO := &SendMessageDTO{
			IsValid:  valid,
			IsSolved: solved,
			Puzzle:   secondPuzzle,
		}

		jsonData, _ := json.Marshal(sendmsgDTO)
		player.Session.Write(jsonData)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		p, ok := s.Get(PlayerString)
		if !ok {
			return
		}
		player := p.(*Player)

		r, ok := s.Get(RoomString)
		if !ok {
			return
		}
		room := r.(*Room)

		room.Mu.Lock()
		if room.Closed {
			room.Mu.Unlock()
			return
		}
		room.Closed = true
		players := room.players
		room.Mu.Unlock()

		for _, p := range players {
			if p != player {
				p.Session.Close()
			}
		}
	})

	r.GET("/ws", func(ctx *gin.Context) {
		m.HandleRequest(ctx.Writer, ctx.Request)
	})

	logger.Println("server is started")
	if err := r.Run(addr); err != nil {
		panic("cannot start server " + err.Error())
	}
}
