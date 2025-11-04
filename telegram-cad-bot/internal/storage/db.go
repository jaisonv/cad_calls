package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps the SQLite database connection
type DB struct {
	conn *sql.DB
}

// NewDB creates a new database connection and initializes the schema
func NewDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// initSchema creates the necessary tables
func (db *DB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		telegram_id INTEGER PRIMARY KEY,
		check_interval INTEGER DEFAULT 5,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS monitored_streets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		telegram_id INTEGER NOT NULL,
		street_name TEXT NOT NULL,
		UNIQUE(telegram_id, street_name),
		FOREIGN KEY(telegram_id) REFERENCES users(telegram_id)
	);

	CREATE TABLE IF NOT EXISTS seen_calls (
		telegram_id INTEGER NOT NULL,
		incident_id TEXT NOT NULL,
		seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY(telegram_id, incident_id),
		FOREIGN KEY(telegram_id) REFERENCES users(telegram_id)
	);

	CREATE INDEX IF NOT EXISTS idx_seen_calls_seen_at ON seen_calls(seen_at);
	`

	_, err := db.conn.Exec(schema)
	return err
}

// GetOrCreateUser gets an existing user or creates a new one
func (db *DB) GetOrCreateUser(telegramID int64) (*User, error) {
	user := &User{}
	err := db.conn.QueryRow(`
		SELECT telegram_id, check_interval, created_at
		FROM users WHERE telegram_id = ?
	`, telegramID).Scan(&user.TelegramID, &user.CheckInterval, &user.CreatedAt)

	if err == sql.ErrNoRows {
		// Create new user
		now := time.Now()
		_, err = db.conn.Exec(`
			INSERT INTO users (telegram_id, check_interval, created_at)
			VALUES (?, 5, ?)
		`, telegramID, now)
		if err != nil {
			return nil, err
		}
		user.TelegramID = telegramID
		user.CheckInterval = 5
		user.CreatedAt = now
		return user, nil
	}

	return user, err
}

// AddMonitoredStreet adds a street to a user's watch list
func (db *DB) AddMonitoredStreet(telegramID int64, streetName string) error {
	// Ensure user exists
	_, err := db.GetOrCreateUser(telegramID)
	if err != nil {
		return err
	}

	// Normalize street name (lowercase for case-insensitive storage)
	streetName = strings.TrimSpace(streetName)

	_, err = db.conn.Exec(`
		INSERT OR IGNORE INTO monitored_streets (telegram_id, street_name)
		VALUES (?, ?)
	`, telegramID, streetName)
	return err
}

// RemoveMonitoredStreet removes a street from a user's watch list
func (db *DB) RemoveMonitoredStreet(telegramID int64, streetName string) error {
	streetName = strings.TrimSpace(streetName)

	result, err := db.conn.Exec(`
		DELETE FROM monitored_streets
		WHERE telegram_id = ? AND street_name = ?
	`, telegramID, streetName)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("street not found in watch list")
	}

	return nil
}

// GetMonitoredStreets returns all streets being monitored by a user
func (db *DB) GetMonitoredStreets(telegramID int64) ([]string, error) {
	rows, err := db.conn.Query(`
		SELECT street_name FROM monitored_streets
		WHERE telegram_id = ?
		ORDER BY street_name
	`, telegramID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var streets []string
	for rows.Next() {
		var street string
		if err := rows.Scan(&street); err != nil {
			return nil, err
		}
		streets = append(streets, street)
	}

	return streets, nil
}

// ClearMonitoredStreets removes all monitored streets for a user
func (db *DB) ClearMonitoredStreets(telegramID int64) error {
	_, err := db.conn.Exec(`
		DELETE FROM monitored_streets WHERE telegram_id = ?
	`, telegramID)
	return err
}

// MarkCallAsSeen marks a call as seen by a user
func (db *DB) MarkCallAsSeen(telegramID int64, incidentID string) error {
	_, err := db.conn.Exec(`
		INSERT OR IGNORE INTO seen_calls (telegram_id, incident_id, seen_at)
		VALUES (?, ?, ?)
	`, telegramID, incidentID, time.Now())
	return err
}

// HasSeenCall checks if a user has seen a specific call
func (db *DB) HasSeenCall(telegramID int64, incidentID string) (bool, error) {
	var count int
	err := db.conn.QueryRow(`
		SELECT COUNT(*) FROM seen_calls
		WHERE telegram_id = ? AND incident_id = ?
	`, telegramID, incidentID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CleanupOldSeenCalls removes seen calls older than the specified number of days
func (db *DB) CleanupOldSeenCalls(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	_, err := db.conn.Exec(`
		DELETE FROM seen_calls WHERE seen_at < ?
	`, cutoff)
	return err
}

// GetAllActiveUsers returns all users who have at least one monitored street
func (db *DB) GetAllActiveUsers() ([]int64, error) {
	rows, err := db.conn.Query(`
		SELECT DISTINCT telegram_id FROM monitored_streets
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []int64
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		users = append(users, userID)
	}

	return users, nil
}

// SetCheckInterval updates the check interval for a user
func (db *DB) SetCheckInterval(telegramID int64, interval int) error {
	// Ensure user exists
	_, err := db.GetOrCreateUser(telegramID)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(`
		UPDATE users SET check_interval = ? WHERE telegram_id = ?
	`, interval, telegramID)
	return err
}
