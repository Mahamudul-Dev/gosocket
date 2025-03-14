package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type Client struct {
	WS       *websocket.Conn
	Username string
	UserID   string
}

type Message struct {
	Type      string `json:"type"`
	UserID    string `json:"user_id"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Target    string `json:"target,omitempty"`
	Timestamp string `json:"timestamp"`
}

type Analytics struct {
	TotalUsers    int `json:"total_users"`
	ActiveUsers   int `json:"active_users"`
	TotalGroups   int `json:"total_groups"`
	TotalMessages int `json:"total_messages"`
}

var (
	clients       = make(map[*Client]bool)
	clientDetails = make(map[string]*Client)
	groupChats    = make(map[string][]*Client)
	messageCount  = 0
	clientMux     sync.Mutex
)

func generateUserID() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func handleClient(ws *websocket.Conn) {
	var username string
	if err := websocket.JSON.Receive(ws, &username); err != nil {
		log.Println("Failed to receive username:", err)
		return
	}

	client := &Client{WS: ws, Username: username, UserID: generateUserID()}
	clientMux.Lock()
	clients[client] = true
	clientDetails[client.UserID] = client
	clientMux.Unlock()

	log.Printf("New client connected: %s (ID: %s)", client.Username, client.UserID)

	for {
		var msg Message
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			log.Println("Error reading message:", err)
			clientMux.Lock()
			delete(clients, client)
			delete(clientDetails, client.UserID)
			clientMux.Unlock()
			return
		}
		msg.Timestamp = time.Now().Format(time.RFC3339)
		handleMessage(client, &msg)
	}
}

func handleMessage(sender *Client, msg *Message) {
	switch msg.Type {
	case "sys-groups":
		listGroups(sender)
	case "sys-group-join":
		groupID := msg.Content
		clientMux.Lock()
		groupChats[groupID] = append(groupChats[groupID], sender)
		clientMux.Unlock()
		sender.WS.Write([]byte(fmt.Sprintf("Joined group %s. Type to chat. Use --sys-exit to leave.\n", groupID)))
	case "sys-exit":

		clientMux.Lock()
		for groupID, members := range groupChats {
			for i, c := range members {
				if c == sender {
					groupChats[groupID] = append(members[:i], members[i+1:]...)
					break
				}
			}
		}
		clientMux.Unlock()
		sender.WS.Close()

	case "sys-myId":
		response := Message{
			Type:    "sys-myId",
			UserID:  sender.UserID,
			Sender:  sender.Username,
			Content: fmt.Sprintf("Your ID: %s, Username: %s", sender.UserID, sender.Username),
		}
		websocket.JSON.Send(sender.WS, response)
	case "sys-peoples":
		listUsers(sender)
	case "sys-analytics":
		getAnalytics(sender)
	case "group":
		groupMessage(sender, msg)
	case "p2p":
		privateMessage(sender, msg)
	default:
		log.Printf("Unknown command: %s", msg.Type)
	}
}

func listGroups(client *Client) {
	clientMux.Lock()
	defer clientMux.Unlock()
	var groupList []string
	for group := range groupChats {
		groupList = append(groupList, group)
	}
	websocket.JSON.Send(client.WS, Message{Type: "sys-groups", Content: fmt.Sprintf("Available groups: %v", groupList)})
}

func privateMessage(sender *Client, msg *Message) {
	clientMux.Lock()
	fmt.Println(clientDetails)
	target, exists := clientDetails[msg.Target]
	fmt.Println(msg.UserID)
	fmt.Println(exists)
	fmt.Println(target)
	clientMux.Unlock()

	if !exists {
		websocket.JSON.Send(sender.WS, Message{
			Type:    "error",
			Content: "User not found!",
		})
		return
	}

	websocket.JSON.Send(sender.WS, msg)
	websocket.JSON.Send(target.WS, msg)
}

func groupMessage(sender *Client, msg *Message) {
	clientMux.Lock()
	members, exists := groupChats[msg.UserID]
	clientMux.Unlock()

	if !exists {
		websocket.JSON.Send(sender.WS, Message{
			Type:    "error",
			Content: "Group does not exist!",
		})
		return
	}

	clientMux.Lock()
	messageCount++
	clientMux.Unlock()

	for _, member := range members {
		if err := websocket.JSON.Send(member.WS, msg); err != nil {
			log.Println("Error sending group message:", err)
		}
	}
}

func getAnalytics(client *Client) {
	clientMux.Lock()
	defer clientMux.Unlock()

	// Prepare the analytics object
	analytics := Analytics{
		TotalMessages: len(groupChats),
		TotalGroups:   len(groupChats),
		TotalUsers:    len(clientDetails),
		ActiveUsers:   len(clients),
	}

	// Marshal the analytics struct into a JSON string
	analyticsJSON, err := json.MarshalIndent(analytics, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling analytics:", err)
		return
	}

	// Print the formatted analytics JSON
	fmt.Println(string(analyticsJSON))

	// Send the formatted JSON response via WebSocket
	websocket.JSON.Send(client.WS, Message{
		Type:    "sys-analytics",
		Content: string(analyticsJSON),
	})
}

func listUsers(client *Client) {
	clientMux.Lock()
	var users []string
	for id, c := range clientDetails {
		users = append(users, fmt.Sprintf("%s (%s)", c.Username, id))
	}
	clientMux.Unlock()

	websocket.JSON.Send(client.WS, Message{
		Type:    "sys-peoples",
		Content: fmt.Sprintf("Active users: %v", users),
	})
}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.Handle("/ws", websocket.Handler(handleClient))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
