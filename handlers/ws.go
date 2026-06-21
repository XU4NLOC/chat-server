package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"chat-server/auth"
	"chat-server/hub"
	"chat-server/config"
	"chat-server/models"
)

var upgrader = websocket.Upgrader {
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(h *hub.Hub, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Validate the token from the query parameters
		tokenStr := r.URL.Query().Get("token")
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(tokenStr, cfg.JWTSecret)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 2. Get the room ID from the query parameters
		roomIDStr := r.URL.Query().Get("room_id")
		if roomIDStr == "" {
			http.Error(w, "Missing room_id", http.StatusBadRequest)
			return
		}

		roomID, err := strconv.Atoi(roomIDStr)
		if err != nil {
			http.Error(w, "Invalid room_id", http.StatusBadRequest)
			return
		}

		// Check if the room exists
		_, err = models.GetRoomByID(roomID)
		if err != nil {
			http.Error(w, "Room not found", http.StatusNotFound)
			return
		}

		// 3. Upgrade the HTTP connection to a WebSocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade to WebSocket: %v", err)
			return
		}

		// 4. Create a new client and register it with the hub
		client := &hub.Client{
			Hub : 		h,
			Conn: 		conn,
			Send: 		make(chan []byte, 256),
			Username: 	claims.Username,
			RoomID: 	roomID,
		}

		h.Register <- client

		// 5. Start the read and write pumps for the client
		go client.WritePump()
		client.ReadPump()
	}
}