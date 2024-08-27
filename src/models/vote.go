package models

import "github.com/google/uuid"

type Vote struct {
	Key      string    `json:"key" db:"key"`
	ClientID string    `json:"client_id" db:"client_id"`
	PollID   uuid.UUID `json:"poll_id" db:"poll_id"`
}
