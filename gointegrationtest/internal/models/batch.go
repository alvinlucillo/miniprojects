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

type Batch struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	DateRequested time.Time          `bson:"name"`
	Status        string             `bson:"status"`
	ErrorMessage  string             `bson:"error_message"`
	Type          string             `bson:"type"`
}

type BatchResponse struct {
	ID            string    `json:"id"`
	DateRequested time.Time `json:"date_requested"`
	Status        string    `json:"status"`
	ErrorMessage  string    `json:"error_message"`
	Type          string    `json:"type"`
}

type BatchRequest struct {
	DateRequested time.Time `json:"date_requested"`
	Type          string    `json:"type"`
}
