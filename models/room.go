package models

import (
	"chat-server/db"
)

type Room struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CreateRoom creates a new room in the database with the given name.
func CreateRoom(name string) (*Room, error) {
	room := &Room{Name: name}
	err := db.DB.QueryRow("INSERT INTO rooms (name) VALUES ($1) RETURNING id", name).Scan(&room.ID)
	return room, err
}

//GetRoomByID retrieves a room from the database by its ID.
func GetRoomByID(id int) (*Room, error) {
	room := &Room{}
	err := db.DB.QueryRow("SELECT id, name FROM rooms WHERE id = $1", id).Scan(&room.ID, &room.Name)
	return room, err
}

//GetAllRooms retrieves all rooms from the database.
func GetAllRooms() ([]Room, error) {
	rows, err := db.DB.Query("SELECT id, name FROM rooms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var r Room
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, err
		}
		rooms = append(rooms, r)
	}
	return rooms, nil
}