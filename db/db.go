package db

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(databaseURL string) {
	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Successfully connected to the database")
	migrate()
}

func migrate() {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
		);

	CREATE TABLE IF NOT EXISTS rooms (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
		);

	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		room_id INT REFERENCES rooms(id),
		user_id INT REFERENCES users(id),
		content TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);
	`
	if _, err := DB.Exec(schema); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}
	log.Println("Database migration completed")
}