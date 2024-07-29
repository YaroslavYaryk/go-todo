package domain

import (
	"time"
)

type Note struct {
	ID        int64     `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Text      string    `json:"text" db:"text"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Tags      string    `json:"tags" db:"tags"`
	Creator   int64     `json:"creator" db:"creator"`
}

type UserNote struct {
	ID     int64 `json:"id" db:"id"`
	UserID int64 `json:"user_id" db:"user_id"`
	NoteID int64 `json:"note_id" db:"note_id"`
}
