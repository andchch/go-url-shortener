package repository

import (
	"auth-service/internal/utils"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// Создаем таблицу при инициализации
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      email TEXT UNIQUE NOT NULL,
      password_hash TEXT NOT NULL,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
  `)
	if err != nil {
		return nil, err
	}
	return &SQLiteRepository{db: db}, nil
}

func (r *SQLiteRepository) CreateUser(email, password string) error {
	hashedPassword := utils.HashPassword(password)
	query := `INSERT INTO users (email, password_hash) VALUES (?, ?)`
	_, err := r.db.Exec(query, email, hashedPassword)
	return err
}

func (r *SQLiteRepository) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, password_hash FROM users WHERE email = ?`
	row := r.db.QueryRow(query, email)
	user := &User{}
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return user, err
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}
