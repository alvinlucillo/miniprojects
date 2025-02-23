package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	BatchStatusPending  = "pending"
	BatchStatusError    = "error"
	BatchStatusComplete = "complete"
)

const (
	BatchTypeGenerateUsersDB = "generate_users_db"
)

type DBExport struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	DateRequested time.Time          `bson:"name"`
	Status        string             `bson:"status"`
	FileName      string             `bson:"file_name"`
	ErrorMessage  string             `bson:"error_message"`
}

type DBExportResponse struct {
	ID            string    `json:"id"`
	DateRequested time.Time `json:"date_requested"`
	FileName      string    `bson:"file_name"`
	Status        string    `json:"status"`
	ErrorMessage  string    `json:"error_message"`
}
