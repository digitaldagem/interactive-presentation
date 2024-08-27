package models

import "github.com/google/uuid"

type Poll struct {
	PollID   uuid.UUID `json:"poll_id"`
	Question string    `json:"question"`
	Options  []Option  `json:"options"`
}

type PollDB struct {
	PollID         uuid.UUID `db:"poll_id"`
	Question       string    `db:"question"`
	PresentationID uuid.UUID `db:"presentation_id"`
	Index          int       `db:"index"`
}
