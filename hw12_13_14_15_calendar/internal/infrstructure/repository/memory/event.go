package memory

import (
	"context"
	"sync"
	"time"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
)

type EventRepository struct {
	mu                    sync.RWMutex
	data                  map[event.ID]event.Model
	participantRepository *EventParticipantRepository
}

var _ event.Repository = (*EventRepository)(nil)

func NewEventRepository(participantRepository *EventParticipantRepository) *EventRepository {
	return &EventRepository{
		data:                  make(map[event.ID]event.Model),
		participantRepository: participantRepository,
	}
}

func (r *EventRepository) Create(_ context.Context, model event.Model) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.data[model.ID]
	if ok {
		return domain.ErrAlreadyExists
	}

	r.data[model.ID] = model
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

	for _, eventID := range eventIds {
		result = append(result, r.data[eventID])
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

	if _, ok := r.data[model.ID]; ok {
		r.data[model.ID] = model
		return nil
	}

	return domain.ErrNotFound
}

func (r *EventRepository) Delete(_ context.Context, model event.Model) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[model.ID]; ok {
		delete(r.data, model.ID)
		return nil
	}

	return domain.ErrNotFound
}

func (r *EventRepository) GetByInterval(ctx context.Context, interval event.UserInterval, excluded ...event.ID) ([]event.Model, error) {
	excludedSet := make(map[event.ID]struct{}, len(excluded))
	for _, excludedID := range excluded {
		excludedSet[excludedID] = struct{}{}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]event.Model, 0)

	eventIds, err := r.participantRepository.GetUserEventIds(ctx, interval.UID)
	if err != nil {
		return nil, err
	}

	for _, eventID := range eventIds {
		if _, ok := excludedSet[eventID]; ok {
			continue
		}

		e := r.data[eventID]
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

func (r *EventRepository) DeleteWhereEndDateBefore(_ context.Context, maxEndDate time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, e := range r.data {
		if e.EndDate.Before(maxEndDate) {
			delete(r.data, id)
		}
	}

	return nil
}

func (r *EventRepository) GetUnprocessedEvents(_ context.Context, interval event.Interval) ([]event.Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []event.Model
	for _, e := range r.data {
		if !e.IsProcessedByScheduler &&
			(e.BeginDate.After(interval.BeginDate) || e.BeginDate.Equal(interval.BeginDate)) &&
			e.BeginDate.Before(interval.EndDate) {
			result = append(result, e)
		}
	}

	return result, nil
}
