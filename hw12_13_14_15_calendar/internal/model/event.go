package model

import (
	"time"

	"github.com/google/uuid"
)

type EventID string

type Event struct {
	ID          EventID   `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	BeginDate   time.Time `db:"begin_date"`
	EndDate     time.Time `db:"end_date"`
}

func NewEvent() Event {
	return Event{
		ID: NewEventID(),
	}
}

func NewEventID() EventID {
	return EventID(
		uuid.New().String(),
	)
}
