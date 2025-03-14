package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type Message struct {
	Type      string `json:"type"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	Target    string `json:"target,omitempty"`
	Timestamp string `json:"timestamp"`
}

func main() {
	serverURL := "ws://localhost:8080/ws"
	ws, err := websocket.Dial(serverURL, "", "http://localhost/")
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}
	defer ws.Close()

	var username string
	fmt.Print("Enter your username: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		username = scanner.Text()
		if err := websocket.JSON.Send(ws, username); err != nil {
			log.Fatal("Failed to send username:", err)
		}
	}

	// Start listening for messages
	go listenForMessages(ws)

	// Read messages from the command line
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			break
		}

		switch text {
		case "--sys-groups":
			msg := Message{
				Type:      "sys-groups",
				Content:   text,
				Timestamp: time.Now().Format(time.RFC3339),
			}

			if err := websocket.JSON.Send(ws, msg); err != nil {
				log.Println("Error sending message:", err)
			}

		case "--sys-peoples":
			msg := Message{
				Type:      "sys-peoples",
				Content:   text,
				Timestamp: time.Now().Format(time.RFC3339),
			}

			if err := websocket.JSON.Send(ws, msg); err != nil {
				log.Println("Error sending message:", err)
			}
		case "--sys-myId":
			msg := Message{
				Type:      "sys-myId",
				Content:   text,
				Timestamp: time.Now().Format(time.RFC3339),
			}

			if err := websocket.JSON.Send(ws, msg); err != nil {
				log.Println("Error sending message:", err)
			}

		case "--sys-analytics":
			msg := Message{
				Type:      "sys-analytics",
				Content:   text,
				Timestamp: time.Now().Format(time.RFC3339),
			}

			if err := websocket.JSON.Send(ws, msg); err != nil {
				log.Println("Error sending message:", err)
			}

		default:
			if strings.HasPrefix(text, "--send-p2p-") {

				parts := strings.Split(text, " ")
				targetId := strings.TrimPrefix(parts[0], "--send-p2p-")
				message := strings.Join(parts[1:], " ")

				fmt.Println("Sending P2P message...")
				fmt.Println("to id:", targetId)
				fmt.Println("message:", message)

				msg := Message{
					Type:      "p2p",
					Target:    targetId,
					Content:   message,
					Timestamp: time.Now().Format(time.RFC3339),
				}

				if err := websocket.JSON.Send(ws, msg); err != nil {
					log.Println("Error sending message:", err)
				}
			} else if strings.HasPrefix(text, "--send-group-") {

				targetId := strings.TrimPrefix(text, "--send-group-")
				msg := Message{
					Type:      "group",
					Target:    targetId,
					Content:   text,
					Timestamp: time.Now().Format(time.RFC3339),
				}

				if err := websocket.JSON.Send(ws, msg); err != nil {
					log.Println("Error sending message:", err)
				}

			} else if strings.HasPrefix(text, "--sys-group-join-") {
				id := strings.TrimPrefix(text, "--sys-group-join-")
				msg := Message{
					Type:      "sys-group-join",
					Content:   id,
					Timestamp: time.Now().Format(time.RFC3339),
				}

				if err := websocket.JSON.Send(ws, msg); err != nil {
					log.Println("Error sending message:", err)
				}
			} else {
				log.Println("Unknown command:", text)
			}
		}

	}
}

func listenForMessages(ws *websocket.Conn) {
	for {
		var msg Message
		fmt.Println("Listening for messages...")
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			log.Println("Connection closed by server")
			break
		}
		fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.Username, msg.Content)
	}
}
