package database

import (
	"database/sql"
	"fmt"
	"gointegrationtest/internal/models"

	_ "modernc.org/sqlite" // SQLite driver
)

func CreateUsersDB(users []models.User, fileName string) error {
	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create users table
	query := `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT
	)`
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	// Insert users into table
	for _, user := range users {
		_, err := db.Exec("INSERT INTO users (id, name) VALUES (?, ?)", user.ID.Hex(), user.Name)
		if err != nil {
			return err
		}
	}

	fmt.Println("SQLite DB created successfully with users")
	return nil
}
