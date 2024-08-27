package models

import "github.com/google/uuid"

type Presentation struct {
	PresentationID   uuid.UUID `json:"presentation_id"`
	CurrentPollIndex int       `json:"current_poll_index"`
	Polls            []Poll    `json:"polls"`
}

type PresentationDB struct {
	PresentationID   uuid.UUID `db:"presentation_id"`
	CurrentPollIndex int       `db:"current_poll_index"`
}
