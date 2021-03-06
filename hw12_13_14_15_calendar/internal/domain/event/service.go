package event

import (
	"context"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
	"github.com/pkg/errors"
)

type Service interface {
	EnsureIntervalAvailable(ctx context.Context, interval UserInterval, excluded ...ID) error
	HasAccess(ctx context.Context, uid user.UID, eventID ID) (bool, error)
	CreateEvent(ctx context.Context, model Model, uids []user.UID) error
	DeleteEvent(ctx context.Context, model Model) error
}

func NewService(eventRepo Repository, participantRepo ParticipantRepository) Service {
	return &service{
		eventRepo:       eventRepo,
		participantRepo: participantRepo,
	}
}

type service struct {
	eventRepo       Repository
	participantRepo ParticipantRepository
}

func (s *service) HasAccess(ctx context.Context, uid user.UID, eventID ID) (bool, error) {
	uids, err := s.participantRepo.GetParticipants(ctx, eventID)
	if err != nil {
		return false, errors.WithStack(err)
	}

	for _, curUID := range uids {
		if uid == curUID {
			return true, nil
		}
	}

	return false, nil
}

func (s *service) EnsureIntervalAvailable(ctx context.Context, interval UserInterval, excluded ...ID) error {
	models, err := s.eventRepo.GetByInterval(ctx, interval, excluded...)
	if err != nil {
		return err
	}

	if len(models) > 0 {
		return ErrIntervalNotAvailable
	}

	return nil
}

func (s *service) CreateEvent(ctx context.Context, model Model, uids []user.UID) error {
	err := s.eventRepo.Create(ctx, model)
	if err != nil {
		return err
	}

	participantModels := make([]Participant, 0, len(uids))
	for _, uid := range uids {
		participantModels = append(participantModels, Participant{
			EventID: model.ID,
			UID:     uid,
		})
	}

	err = s.participantRepo.Create(ctx, participantModels)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteEvent(ctx context.Context, model Model) error {
	err := s.participantRepo.DeleteAllForEvent(ctx, model.ID)
	if err != nil {
		return err
	}

	err = s.eventRepo.Delete(ctx, model)
	if err != nil {
		return err
	}

	return nil
}
