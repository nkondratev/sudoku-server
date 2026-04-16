package main

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestDebug(t *testing.T) {
	player1, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		t.Error("player 1 cannot connect " + err.Error())
	}
	firstMsg := &FirstMessage{}

	if err := player1.ReadJSON(firstMsg); err != nil {
		t.Error(err)
	}

	t.Log(firstMsg)

	player1.Close()
}
