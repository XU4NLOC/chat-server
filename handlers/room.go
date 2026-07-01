package handlers

import (
	"encoding/json"
	"net/http"

	"chat-server/hub"
	"chat-server/models"
)

func CreateRoom(h *hub.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var body struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		room, err := models.CreateRoom(body.Name)
		if err != nil {
			http.Error(w, "could not create room", http.StatusInternalServerError)
			return
		}

		// Broadcast new room to all connected clients
		h.BroadcastGlobal(map[string]interface{}{
			"type":  "system",
			"event": "new_room",
			"room": map[string]interface{}{
				"id":   room.ID,
				"name": room.Name,
			},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(room)
	}
}

func ListRooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rooms, err := models.GetAllRooms()
	if err != nil {
		http.Error(w, "could not fetch rooms", http.StatusInternalServerError)
		return
	}

	if rooms == nil {
		rooms = []models.Room{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

