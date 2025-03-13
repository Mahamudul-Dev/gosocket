package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mahamudul-Dev/gosocket/cmd"
	"golang.org/x/net/websocket"
)

func main() {
	server := cmd.NewServer()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s := websocket.Server{
			Handshake: func(config *websocket.Config, req *http.Request) error {
				// Allow all origins
				config.Origin = req.URL
				return nil
			},
			Handler: websocket.Handler(server.HandleWS),
		}
		s.ServeHTTP(w, r)
	})

	http.HandleFunc("/analytics", func(w http.ResponseWriter, r *http.Request) {
		analytics := server.GetAnalytics()
		json.NewEncoder(w).Encode(analytics)
	})

	fmt.Println("WebSocket server running on ws://localhost:8080")
	http.ListenAndServe(":8080", nil)

}
