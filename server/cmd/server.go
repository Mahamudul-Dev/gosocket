package cmd

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type Client struct {
	conn     *websocket.Conn
	username string
}

type Server struct {
	clients       map[*websocket.Conn]*Client
	groups        map[string]map[*websocket.Conn]bool
	mutex         sync.Mutex
	totalMessages int
}

type Message struct {
	Type      string    `json:"type"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	Target    string    `json:"target"`
	Timestamp time.Time `json:"timestamp"`
}

func NewServer() *Server {
	return &Server{
		clients: make(map[*websocket.Conn]*Client),
		groups:  make(map[string]map[*websocket.Conn]bool),
	}
}

func (s *Server) HandleWS(ws *websocket.Conn) {
	var username string
	websocket.Message.Receive(ws, &username)
	s.mutex.Lock()
	s.clients[ws] = &Client{conn: ws, username: username}
	s.mutex.Unlock()

	fmt.Println(username, "connected")
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	defer func() {
		s.mutex.Lock()
		delete(s.clients, ws)
		s.mutex.Unlock()
		ws.Close()
	}()

	for {
		var msg Message
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		msg.Timestamp = time.Now()
		s.handleMessage(msg, ws)
	}
}

func (s *Server) handleMessage(msg Message, ws *websocket.Conn) {
	s.mutex.Lock()
	s.totalMessages++
	defer s.mutex.Unlock()

	switch msg.Type {
	case "world":
		for clientWS := range s.clients {
			websocket.JSON.Send(clientWS, msg)
		}
	case "group":
		if s.groups[msg.Target] == nil {
			s.groups[msg.Target] = make(map[*websocket.Conn]bool)
		}
		s.groups[msg.Target][ws] = true
		for clientWS := range s.groups[msg.Target] {
			websocket.JSON.Send(clientWS, msg)
		}
	case "private":
		for clientWS, client := range s.clients {
			if client.username == msg.Target {
				websocket.JSON.Send(clientWS, msg)
			}
		}
	}
}

func (s *Server) GetAnalytics() map[string]interface{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return map[string]interface{}{
		"active_users":   len(s.clients),
		"total_groups":   len(s.groups),
		"total_messages": s.totalMessages,
	}
}
