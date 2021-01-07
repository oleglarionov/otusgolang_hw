package event

import (
	"context"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
)

type Service interface {
	EnsureIntervalAvailable(ctx context.Context, interval UserInterval, excluded ...ID) error
	HasAccess(ctx context.Context, uid user.UID, eventId ID) bool
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

func (s *service) HasAccess(ctx context.Context, uid user.UID, eventId ID) bool {
	uids, err := s.participantRepo.GetParticipants(ctx, eventId)
	if err != nil {
		panic(err) // todo?
	}

	for _, curUid := range uids {
		if uid == curUid {
			return true
		}
	}

	return false
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
			EventId: model.Id,
			Uid:     uid,
		})
	}

	err = s.participantRepo.Create(ctx, participantModels)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteEvent(ctx context.Context, model Model) error {
	err := s.participantRepo.DeleteAllForEvent(ctx, model.Id)
	if err != nil {
		return err
	}

	err = s.eventRepo.Delete(ctx, model)
	if err != nil {
		return err
	}

	return nil
}
