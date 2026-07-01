package hub

import (
	"encoding/json"
	"log"

	"chat-server/db"
)

type BroadcastMessage struct {
	RoomID  int
	Sender  string
	Payload []byte
}

type OutgoingMessage struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	RoomID   int    `json:"room_id"`
}

type Hub struct {
	rooms      map[int]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan BroadcastMessage
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[int]map[*Client]bool),
		Broadcast:  make(chan BroadcastMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.rooms[client.RoomID]; !ok {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			log.Printf("Client %s joined room %d", client.Username, client.RoomID)

			h.BroadcastSystemEvent(client.RoomID, map[string]interface{}{
				"type":  "system",
				"event": "user_count",
				"count": len(h.rooms[client.RoomID]),
			})

		case client := <-h.Unregister:
			if clients, ok := h.rooms[client.RoomID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					log.Printf("Client %s left room %d", client.Username, client.RoomID)

					h.BroadcastSystemEvent(client.RoomID, map[string]interface{}{
						"type":  "system",
						"event": "user_count",
						"count": len(h.rooms[client.RoomID]),
					})
				}
			}

		case msg := <-h.Broadcast:

			var incoming struct {
				Content string `json:"content"`
			}

			if err := json.Unmarshal(msg.Payload, &incoming); err != nil {
				log.Printf("Error parsing message from %s: %v", msg.Sender, err)
				continue
			}

			go saveMessage(msg.RoomID, msg.Sender, incoming.Content)

			outgoing, _ := json.Marshal(OutgoingMessage{
				Username: msg.Sender,
				Content:  incoming.Content,
				RoomID:   msg.RoomID,
			})

			for client := range h.rooms[msg.RoomID] {
				select {
				case client.Send <- outgoing:
				default:
					close(client.Send)
					delete(h.rooms[msg.RoomID], client)
				}
			}
		}
	}
}

func saveMessage(roomID int, sender, content string) {
	var userID int
	err := db.DB.QueryRow("SELECT id FROM users WHERE username = $1", sender).Scan(&userID)
	if err != nil {
		log.Printf("Error finding user %s: %v", sender, err)
		return
	}

	_, err = db.DB.Exec("INSERT INTO messages (room_id, user_id, content) VALUES ($1, $2, $3)", roomID, userID, content)
	if err != nil {
		log.Printf("Error saving message from %s: %v", sender, err)
	}
}

func (h *Hub) BroadcastSystemEvent(roomID int, event map[string]interface{}) {
	payload, _ := json.Marshal(event)
	for client := range h.rooms[roomID] {
		select {
		case client.Send <- payload:
		default:
			close(client.Send)
			delete(h.rooms[roomID], client)
		}
	}
}
