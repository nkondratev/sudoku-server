package main

import (
	"encoding/json"
	"log"
	"sudoku-server/sudoku"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

const addr = ":8080"

func main() {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	m := melody.New()
	playerCh := make(chan *Player)

	go func() {

		var room *Room

		for {

			room = NewRoom()
			log.Println("new room")

			for i := range room.players {
				room.players[i] = <-playerCh
			}

			for _, player := range room.players {
				player.Puzzle = sudoku.CopyGrid(room.Puzzle)
			}

			for _, p := range room.players {
				p.Session.Set(room.Id(), room)
			}

			jsonData, err := json.Marshal(room.Puzzle)
			if err != nil {
				for _, player := range room.players {
					player.Session.Close()
				}
				log.Println("game is interrupted")
				continue
			}

			for i := range room.players {
				room.players[i].Session.Write(jsonData)
			}

			log.Println("data is sented")

		}
	}()

	m.HandleConnect(func(s *melody.Session) {
		addr := s.RemoteAddr()
		player := NewPlayer(s)
		log.Printf("new player")
		s.Set(addr.String(), player)
		go func() { playerCh <- player }()
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.Printf("player disconnected")
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {

		log.Println("get message")
		clientMessage := &MessageDTO{}

		if err := json.Unmarshal(msg, clientMessage); err != nil {
			log.Println("cannot read msg")
			return
		}

		p, _ := s.Get(s.Request.RemoteAddr)
		r, _ := s.Get("room")

		player := p.(*Player)
		room := r.(*Room)

		row, col, err := sudoku.ValidAnswer(clientMessage.Puzzle, room.Puzzle)
		if err != nil {
			player.Lives -= 1
		}

		isSolved := sudoku.Equal(player.Puzzle, room.Solution)

		serverMessage := SendMessageDTO{
			Row:      row,
			Col:      col,
			Lives:    player.Lives,
			Puzzle:   player.Puzzle,
			IsSolved: isSolved,
		}

		jsonData, _ := json.Marshal(serverMessage)
		s.Write(jsonData)

		if isSolved {
			log.Printf("player is %v won", player.Session.RemoteAddr())
		}

	})

	r.GET("/ws", func(ctx *gin.Context) {
		m.HandleRequest(ctx.Writer, ctx.Request)
	})

	log.Println("server is started")
	r.Run(addr)
}
