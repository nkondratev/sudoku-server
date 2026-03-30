package main

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestMain(t *testing.T) {

	player1, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		t.Error(err)
	}

	player2, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		t.Error(err)
	}

	firstMsg := &FirstMessage{}

	if err := player1.ReadJSON(firstMsg); err != nil {
		t.Error(err)
	}

	t.Log(firstMsg)

	if err := player2.ReadJSON(firstMsg); err != nil {
		t.Error(err)
	}

	t.Log(firstMsg)

	player1.Close()
	player2.Close()
}
