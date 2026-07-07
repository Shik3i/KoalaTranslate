package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its bcrypt hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSecureToken creates a random hex token.
func GenerateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// CreateSession generates a session token, saves it to database, and returns the token.
func CreateSession(db *sql.DB, userID int, durationHours int) (string, time.Time, error) {
	token, err := GenerateSecureToken()
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := time.Now().Add(time.Duration(durationHours) * time.Hour)

	_, err = db.Exec(
		"INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)",
		token, userID, expiresAt,
	)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

// ValidateSession checks if a session token is valid and returns user info.
func ValidateSession(db *sql.DB, token string) (int, string, error) {
	var userID int
	var username string
	var expiresAt time.Time

	err := db.QueryRow(
		`SELECT s.user_id, u.username, s.expires_at 
		 FROM sessions s 
		 JOIN users u ON s.user_id = u.id 
		 WHERE s.token = ?`,
		token,
	).Scan(&userID, &username, &expiresAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", errors.New("session not found")
		}
		return 0, "", err
	}

	if time.Now().After(expiresAt) {
		// Clean up expired session
		_ = DeleteSession(db, token)
		return 0, "", errors.New("session expired")
	}

	return userID, username, nil
}

// DeleteSession deletes a session from the database.
func DeleteSession(db *sql.DB, token string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE token = ?", token)
	return err
}

// SeedAdmin seeds the admin user in the database if it doesn't already exist.
func SeedAdmin(db *sql.DB, username, password string) error {
	if username == "" || password == "" {
		return errors.New("admin username and password cannot be empty")
	}

	// Check if already seeded
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	if err != nil {
		return err
	}

	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	if exists {
		// Update password in case env variables changed
		_, err = db.Exec("UPDATE users SET password_hash = ?, updated_at = CURRENT_TIMESTAMP WHERE username = ?", hash, username)
		return err
	}

	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, hash)
	return err
}
