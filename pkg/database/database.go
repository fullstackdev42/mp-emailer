package database

import (
	"database/sql"
	"fmt"

	// Import the MySQL driver for database/sql
	_ "github.com/go-sql-driver/mysql"
	"github.com/jonesrussell/loggo"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	*sql.DB
	logger loggo.LoggerInterface
}

func NewDB(dsn string, logger loggo.LoggerInterface, migrationsPath string) (*DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	if err := RunMigrations(dsn, migrationsPath, logger); err != nil {
		return nil, fmt.Errorf("error running migrations: %w", err)
	}

	return &DB{DB: db, logger: logger}, nil
}

func (db *DB) UserExists(username, email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?"
	var count int
	err := db.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %w", err)
	}
	return count > 0, nil
}

func (db *DB) CreateUser(username, email, passwordHash string) error {
	query := "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)"
	_, err := db.Exec(query, username, email, passwordHash)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (db *DB) VerifyUser(username, password string) (bool, error) {
	var storedHash string
	query := "SELECT password_hash FROM users WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // User not found
		}
		return false, fmt.Errorf("error querying user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil // Incorrect password
		}
		return false, fmt.Errorf("error comparing passwords: %w", err)
	}

	return true, nil
}
