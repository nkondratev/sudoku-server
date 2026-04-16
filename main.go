package main

import (
	"app/sudoku"
	"encoding/json"
	"io"
	"log"
	"log/slog"
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

	logger := slog.New(slog.NewJSONHandler(stream, nil))

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	m := melody.New()
	playerCh := make(chan *Player)

	go func() {
		var room *Room
		for {

			room = NewRoom()
			logger.Info("room creates")

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
			logger.Info("puzzle sends to players")
		}
	}()

	m.HandleConnect(func(s *melody.Session) {
		player := NewPlayer(s)
		logger.Info("player connects from client")

		s.Set(PlayerString, player)
		go func() { playerCh <- player }()
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {

		player := s.MustGet(PlayerString).(*Player)
		room := s.MustGet(RoomString).(*Room)

		msgDTO := &MessageDTO{}
		json.Unmarshal(msg, msgDTO)
		logger.Info("message accepts from client")

		player.Mu.Lock()
		player.Puzzle = msgDTO.Puzzle
		puzzleCopy := sudoku.CopyGrid(player.Puzzle)
		player.Mu.Unlock()

		solved := sudoku.IsSolved(puzzleCopy, room.Solution)

		sendmsgDTO := &SendMessageDTO{
			FullPercent: player.Puzzle.FullPercent(room.Solution),
			IsSolved:    solved,
		}
		logger.Info("game finishes")

		jsonData, _ := json.Marshal(sendmsgDTO)
		player.Session.Write(jsonData)
		room.Mu.Lock()
		if room.Closed {
			room.Mu.Unlock()
			logger.Info("room closes")
			return
		}
		room.Closed = true
		players := room.players
		room.Mu.Unlock()

		for _, p := range players {
			p.Session.Close()
		}
		logger.Info("players disconnet")
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
		logger.Info("players disconnet")
	})

	r.GET("/ws", func(ctx *gin.Context) {
		m.HandleRequest(ctx.Writer, ctx.Request)
	})

	logger.Info("server starts")
	if err := r.Run(addr); err != nil {
		logger.Error("failed to start server ")
		panic(err)
	}
}
