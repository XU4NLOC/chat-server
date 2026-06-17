package main

import (
	"log"
	"net/http"

	"chat-server/config"
	"chat-server/db"
    "chat-server/handlers"
)

func main() {
    cfg := config.Load()
    db.Connect(cfg.DatabaseURL)
    mux := http.NewServeMux()
    mux.HandleFunc("/register", handlers.Register(cfg))
    mux.HandleFunc("/login", handlers.Login(cfg))
    
    log.Printf("Server is running on port %s", cfg.Port)
    log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}