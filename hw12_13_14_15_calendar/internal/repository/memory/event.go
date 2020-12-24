package memory

import (
	"context"
	"sync"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/model"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/repository"
)

type EventRepo struct {
	mu   sync.RWMutex
	data map[model.EventID]model.Event
}

var _ repository.EventRepository = (*EventRepo)(nil)

func NewEventRepo() *EventRepo {
	return &EventRepo{
		data: make(map[model.EventID]model.Event),
	}
}

func (r *EventRepo) Create(_ context.Context, event model.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.data[event.ID]
	if ok {
		return repository.ErrAlreadyExists
	}

	r.data[event.ID] = event
	return nil
}

func (r *EventRepo) GetAll(_ context.Context) ([]model.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]model.Event, 0, len(r.data))
	for _, entity := range r.data {
		list = append(list, entity)
	}

	return list, nil
}

func (r *EventRepo) GetByID(_ context.Context, id model.EventID) (model.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	event, ok := r.data[id]
	if !ok {
		return model.Event{}, repository.ErrNotFound
	}

	return event, nil
}

func (r *EventRepo) Update(ctx context.Context, event model.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[event.ID]; ok {
		r.data[event.ID] = event
		return nil
	}

	return repository.ErrNotFound
}

func (r *EventRepo) Delete(_ context.Context, event model.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[event.ID]; ok {
		delete(r.data, event.ID)
		return nil
	}

	return repository.ErrNotFound
}
