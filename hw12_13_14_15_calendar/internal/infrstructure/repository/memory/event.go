package memory

import (
	"context"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
	"sync"
)

type EventRepository struct {
	mu                    sync.RWMutex
	data                  map[event.ID]event.Model
	participantRepository event.ParticipantRepository
}

var _ event.Repository = (*EventRepository)(nil)

func NewEventRepository(participantRepository event.ParticipantRepository) *EventRepository {
	return &EventRepository{
		data:                  make(map[event.ID]event.Model),
		participantRepository: participantRepository,
	}
}

func (r *EventRepository) Create(_ context.Context, model event.Model) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.data[model.Id]
	if ok {
		return domain.ErrAlreadyExists
	}

	r.data[model.Id] = model
	return nil
}

func (r *EventRepository) GetAllForUser(ctx context.Context, uid user.UID) ([]event.Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]event.Model, 0)

	eventIds, err := r.participantRepository.GetUserEventIds(ctx, uid)
	if err != nil {
		return nil, err
	}

	for _, eventId := range eventIds {
		result = append(result, r.data[eventId])
	}

	return result, err
}

func (r *EventRepository) GetByID(_ context.Context, id event.ID) (event.Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	model, ok := r.data[id]
	if !ok {
		return event.Model{}, domain.ErrNotFound
	}

	return model, nil
}

func (r *EventRepository) Update(_ context.Context, model event.Model) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[model.Id]; ok {
		r.data[model.Id] = model
		return nil
	}

	return domain.ErrNotFound
}

func (r *EventRepository) Delete(_ context.Context, model event.Model) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[model.Id]; ok {
		delete(r.data, model.Id)
		return nil
	}

	return domain.ErrNotFound
}

func (r *EventRepository) GetByInterval(ctx context.Context, interval event.UserInterval, excluded ...event.ID) ([]event.Model, error) {
	excludedSet := make(map[event.ID]struct{}, len(excluded))
	for _, excludedId := range excluded {
		excludedSet[excludedId] = struct{}{}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]event.Model, 0)

	eventIds, err := r.participantRepository.GetUserEventIds(ctx, interval.Uid)
	if err != nil {
		return nil, err
	}

	for _, eventId := range eventIds {
		if _, ok := excludedSet[eventId]; ok {
			continue
		}

		e := r.data[eventId]
		if (e.BeginDate.After(interval.BeginDate) || e.BeginDate.Equal(interval.BeginDate)) &&
			e.BeginDate.Before(interval.EndDate) {
			result = append(result, e)
		} else if e.EndDate.After(interval.BeginDate) &&
			(e.EndDate.Before(interval.EndDate) || e.EndDate.Equal(interval.EndDate)) {
			result = append(result, e)
		}
	}

	return result, err
}
