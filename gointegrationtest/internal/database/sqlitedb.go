package database

import (
	"database/sql"
	"fmt"
	"gointegrationtest/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "modernc.org/sqlite" // SQLite driver
)

func CreateUsersDB(users []models.User, fileName string) error {
	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	// Create users table
	query := `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT
	)`
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("creating users table: %w", err)
	}

	// Insert users into table
	for _, user := range users {
		_, err := db.Exec("INSERT INTO users (id, name) VALUES (?, ?)", user.ID.Hex(), user.Name)
		if err != nil {
			return fmt.Errorf("inserting user: %w", err)
		}
	}

	return nil
}

func GetUsersFromDatabase(fileName string) ([]models.User, error) {
	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	// Query all users
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		return nil, fmt.Errorf("querying users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var id string
		if err := rows.Scan(&id, &user.Name); err != nil {
			return nil, fmt.Errorf("scanning user: %w", err)
		}
		user.ID, err = primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("converting id [%s] to ObjectID: %w", id, err)
		}
		users = append(users, user)
	}

	return users, nil
}
