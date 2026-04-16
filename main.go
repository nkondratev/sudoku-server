package main

import (
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

	server := &Server{
		m:        m,
		playerCh: make(chan *Player, 100),
		logger:   logger,
	}

	go server.matchmakingLoop()

	m.HandleConnect(func(s *melody.Session) {
		player := NewPlayer(s)
		logger.Info("player connected")

		s.Set(PlayerString, player)

		// неблокирующая отправка
		select {
		case server.playerCh <- player:
		default:
			logger.Warn("player channel full")
			_ = s.Close()
		}
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		playerAny, ok1 := s.Get(PlayerString)
		roomAny, ok2 := s.Get(RoomString)

		if !ok1 || !ok2 {
			return
		}

		room := roomAny.(*Room)

		room.Mu.Lock()
		if room.IsEnd {
			room.Mu.Unlock()
			return
		}

		room.IsEnd = true
		room.Mu.Unlock()

		player1 := playerAny.(*Player)
		var player2 *Player
		for _, p := range room.players {
			if player1.Session != p.Session {
				player2 = p
				break
			}
		}

		var msgDTO MessageDTO
		if err := json.Unmarshal(msg, &msgDTO); err != nil {
			logger.Error("cannot get MessageDTO")
			return
		}

		if player2 == nil {
			player1.Session.Write(func() []byte {
				data, _ := json.Marshal(MESSAGE{
					Result: "WIN",
				})
				return data
			}(),
			)
			room.Close(logger)
			return
		}

		player1.Mu.Lock()
		player1.Puzzle = msgDTO.Puzzle
		player1.Mu.Unlock()

		percent1 := player1.Puzzle.FullPercent(room.Solution)
		percent2 := player2.Puzzle.FullPercent(room.Solution)

		switch {
		case percent1 > percent2:
			player1.Session.Write(func() []byte {
				data, _ := json.Marshal(MESSAGE{
					Result: "WIN",
				})
				return data
			}())
			player2.Session.Write(func() []byte {
				data, _ := json.Marshal(MESSAGE{
					Result: "LOSE",
				})
				return data
			}())
		case percent1 < percent2:
			player1.Session.Write(func() []byte {
				data, _ := json.Marshal(MESSAGE{
					Result: "LOSE",
				})
				return data
			}())
			player2.Session.Write(func() []byte {
				data, _ := json.Marshal(MESSAGE{
					Result: "WIN",
				})
				return data
			}())
		case percent1 == percent2:
			player1.Session.Write(func() []byte {
				data, _ := json.Marshal(MESSAGE{
					Result: "DRAW",
				})
				return data
			}())
			player2.Session.Write(func() []byte {
				data, _ := json.Marshal(MESSAGE{
					Result: "DRAW",
				})
				return data
			}())

		}

		logger.Info("message processed")

		room.Close(logger)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		pAny, ok := s.Get(PlayerString)
		if !ok {
			return
		}
		player := pAny.(*Player)

		roomAny, ok := s.Get(RoomString)
		if !ok {
			return
		}
		room := roomAny.(*Room)

		logger.Info("player disconnected")

		room.Close(logger)

		_ = player
	})

	r.GET("/ws", func(ctx *gin.Context) {
		m.HandleRequest(ctx.Writer, ctx.Request)
	})

	logger.Info("server started")

	if err := r.Run(addr); err != nil {
		logger.Error("server failed", "err", err)
		panic(err)
	}
}

// MATCHMAKING LOOP
