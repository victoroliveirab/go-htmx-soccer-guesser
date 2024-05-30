package models

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"strings"
)

type User struct {
	Id           int
	Username     string
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    string
	UpdatedAt    string
}

// Unsafe: should use bcrypt in real world app
func generateSalt(length int) (string, error) {
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func hashPassword(password, salt string) string {
	saltedPassword := password + salt
	hash := sha256.Sum256([]byte(saltedPassword))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func GetLoggingInUser(db *sql.DB, username, password string) *User {
	var user User
	row := db.QueryRow("SELECT * FROM Users WHERE username = $1", username)
	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil
	}
	passwordHashSplit := strings.Split(user.PasswordHash, ":")
	salt := passwordHashSplit[0]
	storedHash := passwordHashSplit[1]

	hashedIncomingPassword := hashPassword(password, salt)
	if storedHash != hashedIncomingPassword {
		return nil
	}
	return &user
}

func CreateUser(db *sql.DB, username, email, password string) (int64, error) {
	salt, err := generateSalt(16)
	if err != nil {
		return -1, err
	}
	passwordHash := hashPassword(password, salt)
	saltedPasswordHash := salt + ":" + passwordHash
	row, err := db.Exec(
		"INSERT INTO Users(username, email, password_hash) VALUES ($1, $2, $3)",
		username,
		email,
		saltedPasswordHash,
	)
	if err != nil {
		return -1, err
	}
	return row.LastInsertId()
}

func GetUserById(db *sql.DB, id int64) (*User, error) {
	row := db.QueryRow("SELECT * FROM Users WHERE id = $1", id)

	var user User
	if err := row.Scan(&user.Id, &user.Username, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}
