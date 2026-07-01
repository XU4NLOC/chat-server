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
	all        map[*Client]bool // every connected client, across all rooms
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan BroadcastMessage
	Global     chan []byte // broadcasts to ALL clients regardless of room
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[int]map[*Client]bool),
		all:        make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan BroadcastMessage),
		Global:     make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			// Add to room set
			if _, ok := h.rooms[client.RoomID]; !ok {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true

			// Add to global set
			h.all[client] = true
			log.Printf("✅ %s joined room %d", client.Username, client.RoomID)

			h.broadcastSystemEvent(client.RoomID, map[string]interface{}{
				"type":  "system",
				"event": "user_count",
				"count": len(h.rooms[client.RoomID]),
			})

		case client := <-h.Unregister:
			if clients, ok := h.rooms[client.RoomID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					delete(h.all, client) // remove from global set too
					close(client.Send)
					log.Printf("❌ %s left room %d", client.Username, client.RoomID)

					h.broadcastSystemEvent(client.RoomID, map[string]interface{}{
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
				log.Printf("Bad message from %s: %v", msg.Sender, err)
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
					delete(h.all, client)
				}
			}

		case payload := <-h.Global:
			// Send to every connected client across all rooms
			for client := range h.all {
				select {
				case client.Send <- payload:
				default:
					close(client.Send)
					delete(h.all, client)
				}
			}
		}
	}
}

func (h *Hub) broadcastSystemEvent(roomID int, event map[string]interface{}) {
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

func (h *Hub) BroadcastGlobal(event map[string]interface{}) {
	payload, _ := json.Marshal(event)
	h.Global <- payload
}

func saveMessage(roomID int, username, content string) {
	var userID int
	err := db.DB.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		log.Printf("saveMessage: user not found: %v", err)
		return
	}
	_, err = db.DB.Exec(
		"INSERT INTO messages (room_id, user_id, content) VALUES ($1, $2, $3)",
		roomID, userID, content,
	)
	if err != nil {
		log.Printf("saveMessage: insert failed: %v", err)
	}
}
