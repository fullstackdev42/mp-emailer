package models

import "time"

type Campaign struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Template  string    `json:"template"`
	OwnerID   int       `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tokens    []string  `json:"tokens"`
}
