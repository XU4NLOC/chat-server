package main

import (
	"log"
	"net/http"

	"chat-server/config"
	"chat-server/db"
	"chat-server/handlers"
	"chat-server/hub"
)

func main() {
	cfg := config.Load()
	db.Connect(cfg.DatabaseURL)

	h := hub.NewHub()
	go h.Run()

	mux := http.NewServeMux()

	mux.HandleFunc("/register", handlers.Register(cfg))
	mux.HandleFunc("/login", handlers.Login(cfg))

	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateRoom(h)(w, r) // pass hub in
		case http.MethodGet:
			handlers.ListRooms(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/ws", handlers.ServeWS(h, cfg))
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	log.Printf("Server running on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}

