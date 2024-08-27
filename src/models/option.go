package models

import "github.com/google/uuid"

type Option struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OptionDB struct {
	Key    string    `db:"key"`
	Value  string    `db:"value"`
	PollID uuid.UUID `db:"poll_id"`
	Index  int       `db:"index"`
}
