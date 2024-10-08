package models

import "time"

type Campaign struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Template  string    `json:"template"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
