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
	roomCh := make(chan *Room)

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

			roomCh <- room

		}
	}()

	m.HandleConnect(func(s *melody.Session) {
		player := NewPlayer(s)
		log.Printf("new player")
		go func() { playerCh <- player }()
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.Printf("player disconnected")
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		log.Println("get message")
		clientMessage := &SendMesageDTO{}

		if err := json.Unmarshal(msg, clientMessage); err != nil {
			log.Println("cannot read msg")
		}

		room := <-roomCh

		//TODO find player
		row, col, err := sudoku.ValidAnswer(clientMessage.Puzzle, room.Puzzle)
		if err != nil {
		}
		serverMessage := GetMessageDTO{
			Row: row,
			Col: col,
		}
		for _, p := range room.players {
			if p.Session == s {
				p.Lives -= 1
			}
		}

		s.Write()
		//TODO
	})

	r.GET("/ws", func(ctx *gin.Context) {
		m.HandleRequest(ctx.Writer, ctx.Request)
	})

	r.POST("/finish", func(ctx *gin.Context) {
		log.Println("game is ended")
		//TODO
	})

	log.Println("server is started")

	r.Run(addr)
}
