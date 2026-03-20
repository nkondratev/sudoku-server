package main

import (
	"encoding/json"
	"log"
	"net/http"
	ws "sudoku-server/webSocket"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

const (
	addr = ":8080"
)

func main() {
	r := gin.Default()
	server := ws.NewServer()
	m := melody.New()
	var (
		matchmaking = make(chan *ws.Player)
		game        = make(chan *ws.Room)
	)

	go ws.HandleRooms(matchmaking, game)

	go func() {
		room := <-game
		server.Append(room)
		ws.StartGame(room)
	}()

	r.GET("/ws", func(ctx *gin.Context) {
		m.HandleRequest(ctx.Writer, ctx.Request)
	})

	m.HandleConnect(func(s *melody.Session) {
		player := ws.NewPlayer(s)
		matchmaking <- player
	})

	m.HandleDisconnect(func(s *melody.Session) {
		server.DeletePlayer(s)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		moveDTO := &ws.MoveDTO{}
		if err := json.Unmarshal(msg, moveDTO); err != nil {
			log.Printf("%v", err)
			return
		}
		server.Players[s].Puzzle.Grid = moveDTO.Sudoku.Grid
		room, err := server.FindRoom(s)
		if err != nil {
			log.Printf("%v", err)
			return
		}
	})

	r.POST("/finish", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": true,
		})
	})

	r.Run(addr)
}
