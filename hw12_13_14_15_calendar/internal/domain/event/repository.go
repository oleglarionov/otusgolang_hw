package event

import (
	"context"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
)

type Repository interface {
	Create(ctx context.Context, model Model) error
	GetByID(ctx context.Context, id ID) (Model, error)
	Update(ctx context.Context, model Model) error
	Delete(ctx context.Context, model Model) error
	GetByInterval(ctx context.Context, interval UserInterval, excluded ...ID) ([]Model, error)
}

type ParticipantRepository interface {
	Create(ctx context.Context, participants []Participant) error
	DeleteAllForEvent(ctx context.Context, eventID ID) error
	GetUserEventIds(ctx context.Context, uid user.UID) ([]ID, error)
	GetParticipants(ctx context.Context, eventID ID) ([]user.UID, error)
}
