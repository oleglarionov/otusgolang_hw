package sql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
)

type EventParticipantRepository struct {
	db  *sqlx.DB
	sbt sq.StatementBuilderType
}

func NewEventParticipantRepository(db *sqlx.DB, sbt sq.StatementBuilderType) event.ParticipantRepository {
	return &EventParticipantRepository{db: db, sbt: sbt}
}

func (r *EventParticipantRepository) Create(ctx context.Context, participants []event.Participant) error {
	if len(participants) == 0 {
		return nil
	}

	qb := r.sbt.Insert("event_participants").
		Columns("event_id", "uid")
	for _, p := range participants {
		qb = qb.Values(p.EventID, p.UID)
	}
	qb = qb.Suffix("on conflict do nothing")

	_, err := qb.RunWith(r.db).ExecContext(ctx)
	return err
}

func (r *EventParticipantRepository) GetUserEventIds(ctx context.Context, uid user.UID) ([]event.ID, error) {
	panic("implement me")
}

func (r *EventParticipantRepository) GetParticipants(ctx context.Context, eventID event.ID) ([]user.UID, error) {
	query, args := r.sbt.Select("uid").
		From("event_participants").
		Where(sq.Eq{"event_id": eventID}).
		MustSql()

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]user.UID, 0)
	for rows.Next() {
		var uid user.UID
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		result = append(result, uid)
	}

	return result, nil
}

func (r *EventParticipantRepository) DeleteAllForEvent(ctx context.Context, eventID event.ID) error {
	_, err := r.db.ExecContext(ctx, "delete from event_participants where event_id = $1", &eventID)
	if err != nil {
		return fmt.Errorf("error deleting event participants: %w", err)
	}

	return nil
}
