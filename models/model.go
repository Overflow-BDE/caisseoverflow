package models

import "time"

// Model is the standard model for what's stored in DB
type Model struct {
	ID        int64     `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
