package repository

import (
	"auth-service/internal/utils"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(email, password string) error {
	hashedPassword := utils.HashPassword(password)
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2)`
	_, err := r.db.Exec(query, email, hashedPassword)
	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, password_hash FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)
	user := &User{}
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return user, err
}
