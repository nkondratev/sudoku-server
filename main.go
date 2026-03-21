package main

import (
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"

	ws "sudoku-server/websocket"
)

const (
	addr = ":8080"
)

func main() {
	r := gin.Default()
	m := melody.New()
	server := ws.NewServer()

	r.GET("/ws", func(ctx *gin.Context) {
		m.HandleRequest(ctx.Writer, ctx.Request)
	})

	m.HandleConnect(func(s *melody.Session) {
		room := server.FindOrCreateRoom(s)
		room.StartGame()
	})

	m.HandleDisconnect(func(s *melody.Session) {
		server.DisconnectPlayer(s)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
	})

	r.GET("/finish", func(ctx *gin.Context) {
	})

	r.Run(addr)
}
