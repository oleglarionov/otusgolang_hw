package scheduler

import (
	"context"
	"time"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/application/broker"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/notification"
	"github.com/pkg/errors"
)

type Notifier interface {
	Notify(ctx context.Context) error
}

var _ Notifier = (*NotifierImpl)(nil)

type NotifierImpl struct {
	pusher          broker.Pusher
	eventRepo       event.Repository
	participantRepo event.ParticipantRepository
}

func NewNotifierImpl(
	pusher broker.Pusher,
	eventRepo event.Repository,
	participantRepo event.ParticipantRepository,
) *NotifierImpl {
	return &NotifierImpl{
		pusher:          pusher,
		eventRepo:       eventRepo,
		participantRepo: participantRepo,
	}
}

func (n *NotifierImpl) Notify(ctx context.Context) error {
	events, err := n.eventRepo.GetUnprocessedEvents(ctx, event.Interval{
		BeginDate: time.Now(),
		EndDate:   time.Now().Add(time.Hour),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	for _, e := range events {
		participants, err := n.participantRepo.GetParticipants(ctx, e.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		ntfcE := notification.Event{
			Title:       e.Title,
			Description: e.Description,
			BeginDate:   e.BeginDate,
			EndDate:     e.EndDate,
		}

		for _, p := range participants {
			ntfc := notification.Model{
				UID:   p,
				Event: ntfcE,
			}

			err := n.pusher.Push(ctx, ntfc)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		e.IsProcessedByScheduler = true
		err = n.eventRepo.Update(ctx, e)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
