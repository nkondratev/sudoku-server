package main

import (
	"net/http"
	ws "server/webSocket"

	"github.com/olahol/melody"
)

func main() {
	var server = ws.NewServer()
	var matchmaking = make(chan *ws.Player)
	m := melody.New()

	go func() {
		for {
			ws.HandleRooms(matchmaking, server)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})

	m.HandleConnect(func(s *melody.Session) {
		player := ws.NewPlayer(s)
		server.AppendPlayer(player)
		matchmaking <- player
	})

	m.HandleDisconnect(func(s *melody.Session) {
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
	})

	http.ListenAndServe(":5000", nil)
}
