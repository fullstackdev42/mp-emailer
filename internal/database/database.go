package database

import (
	"database/sql"
	"errors"
	"fmt"

	// Import the MySQL driver for database/sql
	_ "github.com/go-sql-driver/mysql"
	"github.com/jonesrussell/loggo"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	SQL    *sql.DB
	logger loggo.LoggerInterface
}

func NewDB(dsn string, logger loggo.LoggerInterface) (*DB, error) {
	logger.Debug("Attempting to connect to database with DSN: " + dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		logger.Error("error connecting to database", err)
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	logger.Info("Successfully connected to database")

	return &DB{SQL: db, logger: logger}, nil
}

func (db *DB) UserExists(username, email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?"
	var count int
	err := db.SQL.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %w", err)
	}
	return count > 0, nil
}

func (db *DB) CreateUser(id, username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	query := "INSERT INTO users (id, username, email, password_hash) VALUES (?, ?, ?, ?)"
	_, err = db.SQL.Exec(query, id, username, email, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (db *DB) VerifyUser(username, password string) (string, error) {
	var storedHash, userID string
	query := "SELECT id, password_hash FROM users WHERE username = ?"
	err := db.SQL.QueryRow(query, username).Scan(&userID, &storedHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("invalid username or password")
		}
		return "", fmt.Errorf("error querying user: %w", err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid username or password")
	}
	return userID, nil
}
