package models

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"chat-server/db"
)

type User struct {
	ID            int    
	Username      string 
	PasswordHash  string
}

// CreateUser creates a new user in the database with the given username and password.
func CreateUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	_, err = db.DB.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, string(hash))
	return err
}

// GetUserByUsername retrieves a user from the database by their username.
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := db.DB.QueryRow("SELECT id, username, password_hash FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != sql.ErrNoRows {
		return user, err
	}
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// CheckPassword checks if the provided password matches the user's password hash.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
