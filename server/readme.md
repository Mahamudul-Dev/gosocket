# CLI Chat Server Documentation

## Overview
This project is a CLI-based chat server built using Go and WebSockets. It supports:
- **World Chat** (messages visible to all users)
- **Private Messaging** (one-to-one chat)
- **Group Chat** (multiple users in a chatroom)
- **Analytics** (track active users, total groups, and total messages sent)

## How It Works
1. The server listens on `ws://localhost:8080/ws` for WebSocket connections.
2. Clients connect, provide a username, and can send messages to the world, a specific user, or a group.
3. The server processes messages and routes them accordingly.
4. An analytics endpoint is available at `http://localhost:8080/analytics`.

---

## Running the Server
### Prerequisites
Ensure you have Go installed.

### Steps
1. Clone the repository:
   ```sh
   git clone https://github.com/Mahamudul-Dev/gosocket.git
   cd gosocket
   ```
2. Run the server:
   ```sh
   go run main.go
   ```
3. The WebSocket server starts on `ws://localhost:8080/ws`.
4. Check analytics at `http://localhost:8080/analytics`.

---

## Connecting a Client
You can use WebSocket clients like `wscat` (Node.js package) or a Go-based client.

### Using `wscat`
1. Install wscat (if not installed):
   ```sh
   npm install -g wscat
   ```
2. Connect to the server:
   ```sh
   wscat -c ws://localhost:8080/ws
   ```
3. Provide a username and start chatting!

---

## Message Format
Messages should be JSON-formatted:
```json
{
  "type": "world | private | group",
  "sender": "your_username",
  "content": "Hello, World!",
  "target": "target_username_or_group",
  "timestamp": "auto-generated"
}
```

### Message Types
- **World Chat**: Set `type` to `world`.
- **Private Chat**: Set `type` to `private` and specify `target` as the recipient.
- **Group Chat**: Set `type` to `group` and specify `target` as the group name.

---

## Analytics Endpoint
Retrieve chat analytics:
```sh
curl http://localhost:8080/analytics
```
Response example:
```json
{
  "active_users": 5,
  "total_groups": 3,
  "total_messages": 120
}
```

---

## Future Enhancements
- Authentication for users
- Persistent storage for chat history
- Advanced group management

Happy Coding! 🚀

