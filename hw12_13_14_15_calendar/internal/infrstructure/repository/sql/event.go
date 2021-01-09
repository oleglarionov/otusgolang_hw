package sql

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
)

type EventRepository struct {
	db  *sqlx.DB
	sbt sq.StatementBuilderType
}

var _ event.Repository = (*EventRepository)(nil)

func NewEventRepository(db *sqlx.DB, sbt sq.StatementBuilderType) event.Repository {
	return &EventRepository{db: db, sbt: sbt}
}

func (r *EventRepository) Create(ctx context.Context, model event.Model) error {
	_, err := r.sbt.Insert("events").
		Columns("id", "title", "description", "begin_date", "end_date").
		Values(model.Id, model.Title, model.Description, model.BeginDate, model.EndDate).
		RunWith(r.db).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *EventRepository) GetByID(ctx context.Context, id event.ID) (event.Model, error) {
	row := r.db.QueryRowxContext(ctx, "select * from events where id = $1", id)
	model := event.Model{}
	if err := row.StructScan(&model); err != nil {
		return model, fmt.Errorf("error getting event: %w", err)
	}

	return model, nil
}

func (r *EventRepository) Update(ctx context.Context, model event.Model) error {
	_, err := r.db.NamedExecContext(ctx, "update events "+
		"set title = :title, "+
		"description = :description, "+
		"begin_date = :begin_date, "+
		"end_date = :end_date "+
		"where id = :id",
		&model,
	)
	if err != nil {
		return fmt.Errorf("error updating event: %w", err)
	}

	return nil
}

func (r *EventRepository) Delete(ctx context.Context, model event.Model) error {
	_, err := r.db.NamedExecContext(ctx, "delete from events where id = :id", &model)
	if err != nil {
		return fmt.Errorf("error deleting event: %w", err)
	}

	return nil
}

func (r *EventRepository) GetByInterval(ctx context.Context, interval event.UserInterval, excluded ...event.ID) ([]event.Model, error) {
	qb := sq.Select("events.*").
		From("events").
		InnerJoin("event_participants "+
			"on events.id = event_participants.event_id "+
			"and event_participants.uid = ?",
			interval.Uid).
		Where(
			sq.Or{
				sq.And{
					sq.GtOrEq{"begin_date": interval.BeginDate},
					sq.Lt{"begin_date": interval.EndDate},
				},
				sq.And{
					sq.Gt{"end_date": interval.BeginDate},
					sq.LtOrEq{"end_date": interval.EndDate},
				},
			})
	if len(excluded) > 0 {
		qb = qb.Where(sq.Expr("id not in (?)", excluded))
	}
	query, args := qb.MustSql()

	query, args, err := sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting models: %w", err)
	}
	defer rows.Close()

	models := make([]event.Model, 0)
	for rows.Next() {
		model := event.Model{}
		if err := rows.StructScan(&model); err != nil {
			return nil, fmt.Errorf("error scanning event: %w", err)
		}

		models = append(models, model)
	}

	return models, nil
}
