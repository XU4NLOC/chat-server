package models

import "chat-server/db"

type Message struct {
	Username  string `json:"username"`
	Content   string `json:"content"`
	RoomID    int    `json:"room_id"`
	CreatedAt string `json:"created_at"`
}

func GetRecentMessages(roomID, limit int) ([]Message, error) {
	rows, err := db.DB.Query(`
		SELECT u.username, m.Content, m.room_id, m.created_at
		FROM messages m
		JOIN users u ON u.id = m.user_id
		WHERE m.room_id = $1
		ORDER BY m.created_at DESC
		LIMIT $2
		`, roomID, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.Username, &msg.Content, &msg.RoomID, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

