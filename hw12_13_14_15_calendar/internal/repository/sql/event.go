package sql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/model"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/repository"
)

type EventRepo struct {
	db *sqlx.DB
}

var _ repository.EventRepository = (*EventRepo)(nil)

func NewEventRepository(db *sqlx.DB) *EventRepo {
	return &EventRepo{db: db}
}

func (r *EventRepo) Create(ctx context.Context, event model.Event) error {
	_, err := r.db.ExecContext(ctx, "insert into events "+
		"(id, title, description, begin_date, end_date) "+
		"values (:id, :title, :description, :begin_date, :end_date)",
		event,
	)
	if err != nil {
		return fmt.Errorf("error creating event: %w", err)
	}

	return nil
}

func (r *EventRepo) GetAll(ctx context.Context) ([]model.Event, error) {
	rows, err := r.db.QueryxContext(ctx, "select * from events order by begin_date desc")
	if err != nil {
		return nil, fmt.Errorf("error getting events: %w", err)
	}
	defer rows.Close()

	events := make([]model.Event, 0)
	for rows.Next() {
		event := model.Event{}
		if err := rows.StructScan(&event); err != nil {
			return nil, fmt.Errorf("error scanning event: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepo) GetByID(ctx context.Context, id model.EventID) (model.Event, error) {
	row := r.db.QueryRowxContext(ctx, "select * from events where id = ?", id)
	event := model.Event{}
	if err := row.StructScan(&event); err != nil {
		return event, fmt.Errorf("error getting event: %w", err)
	}

	return event, nil
}

func (r *EventRepo) Update(ctx context.Context, event model.Event) error {
	_, err := r.db.ExecContext(ctx, "update events "+
		"set title = :title, "+
		"description = :description "+
		"begin_date = :begin_date "+
		"end_date = :end_date "+
		"where id = :id",
		&event,
	)
	if err != nil {
		return fmt.Errorf("error updating event: %w", err)
	}

	return nil
}

func (r *EventRepo) Delete(ctx context.Context, event model.Event) error {
	_, err := r.db.ExecContext(ctx, "delete from events where id = :id", &event)
	if err != nil {
		return fmt.Errorf("error deleting event: %w", err)
	}

	return nil
}
