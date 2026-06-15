package main

import (
    "log"

    "chat-server/config"
    "chat-server/db"
)

func main() {
    cfg := config.Load()
    db.Connect(cfg.DatabaseURL)

    log.Printf("Server will run on port %s", cfg.Port)
}