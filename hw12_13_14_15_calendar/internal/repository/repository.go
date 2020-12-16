package repository

import (
	"context"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/model"
)

type EventRepository interface {
	Create(ctx context.Context, event model.Event) error
	GetAll(ctx context.Context) ([]model.Event, error)
	GetByID(ctx context.Context, id model.EventID) (model.Event, error)
	Update(ctx context.Context, event model.Event) error
	Delete(ctx context.Context, event model.Event) error
}

type Repository struct {
	EventRepository
}
