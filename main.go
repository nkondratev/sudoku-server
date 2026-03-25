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

	var hub = make([]*Room, 0, 1024)

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

			hub = append(hub, room)
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

		go func() {

			log.Println("get message")
			clientMessage := &MessageDTO{}

			if err := json.Unmarshal(msg, clientMessage); err != nil {
				log.Println("cannot read msg")
				return
			}

			var room *Room
			var player *Player

			var isFind bool = false
			for _, r := range hub {
				if !isFind {
					for _, p := range r.players {
						if p.Session == s {
							room = r
							player = p
							isFind = true
							break
						}
					}
				} else {
					break
				}
			}

			row, col, err := sudoku.ValidAnswer(clientMessage.Puzzle, room.Puzzle)
			if err != nil {
				player.Lives -= 1
			}

			serverMessage := SendMessageDTO{
				Row:      row,
				Col:      col,
				Lives:    player.Lives,
				Puzzle:   player.Puzzle,
				IsSolved: sudoku.Equal(player.Puzzle, room.Solution),
			}

			jsonData, err := json.Marshal(serverMessage)
			if err != nil {
				//TODO
			}

			s.Write(jsonData)

		}()

	})

	r.GET("/ws", func(ctx *gin.Context) {
		m.HandleRequest(ctx.Writer, ctx.Request)
	})

	log.Println("server is started")
	r.Run(addr)
}
