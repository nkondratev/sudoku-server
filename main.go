package main

import (
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

	r.GET("/ws", func(ctx *gin.Context) {
		m.HandleRequest(ctx.Writer, ctx.Request)
	})

	m.HandleConnect(func(s *melody.Session) {
	})

	m.HandleDisconnect(func(s *melody.Session) {
		server.DeletePlayer(s)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
	})

	r.POST("/finish", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": true,
		})
	})

	r.Run(addr)
}
