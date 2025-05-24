package main

import (
	"echo/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading: ", err)
		return
	}

	defer conn.Close()
	var cfg utils.Config

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Cannot read message: ", err)
			return
		}

		fmt.Printf("Received message: %s\n", message)
		if err := json.Unmarshal(message, &cfg); err != nil {
			fmt.Println("Cannot parse json")
			return
		}

		localAddr := fmt.Sprintf(":%s", cfg.LocalPort) // TODO Add support for receiver
		if err := RunPeer(localAddr, cfg.RemoteAddr, cfg.FilePath, cfg.Dest, false); err != nil {
			fmt.Println("Cannot run peer: ", err)
			return
		}

		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error writing message: ", err)
			return
		}
	}
}
