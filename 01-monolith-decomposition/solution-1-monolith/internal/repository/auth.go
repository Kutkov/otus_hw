package repository

import (
	"database/sql"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(database *sql.DB) *AuthRepository {
	return &AuthRepository{db: database}
}

func (r *AuthRepository) CreateToken(token, userID, createdAt string) error {
	_, err := r.db.Exec(`INSERT INTO tokens(token, user_id, created_at) VALUES(?,?,?)`, token, userID, createdAt)
	return err
}

func (r *AuthRepository) GetUserIDByToken(token string) (string, error) {
	var userID string
	err := r.db.QueryRow("SELECT user_id FROM tokens WHERE token = ?", token).Scan(&userID)
	return userID, err
}

func (r *AuthRepository) GetPasswordHash(userID string) ([]byte, error) {
	var passwordHash []byte
	err := r.db.QueryRow("SELECT password_hash FROM users WHERE id = ?", userID).Scan(&passwordHash)
	return passwordHash, err
}
