package repository

import (
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(database *sql.DB) *UserRepository {
	return &UserRepository{db: database}
}

type User struct {
	ID           string
	FirstName    string
	SecondName   string
	Birthdate    string
	Biography    string
	City         string
	PasswordHash []byte
	CreatedAt    string
}

func (r *UserRepository) Create(user *User) error {
	_, err := r.db.Exec(
		`INSERT INTO users(id, first_name, second_name, birthdate, biography, city, password_hash, created_at) VALUES(?,?,?,?,?,?,?,?)`,
		user.ID, user.FirstName, user.SecondName, user.Birthdate, user.Biography, user.City, user.PasswordHash, user.CreatedAt,
	)
	return err
}

func (r *UserRepository) GetByID(id string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(
		`SELECT id, first_name, second_name, birthdate, biography, city, password_hash, created_at FROM users WHERE id = ?`,
		id,
	).Scan(&user.ID, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Exists(id string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT 1 FROM users WHERE id = ?", id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return exists, err
}
