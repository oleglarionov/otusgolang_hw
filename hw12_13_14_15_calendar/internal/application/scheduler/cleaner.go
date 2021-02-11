package scheduler

import (
	"context"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/pkg/errors"
	"time"
)

type Cleaner interface {
	Clean(ctx context.Context) error
}

var _ Cleaner = (*CleanerImpl)(nil)

type CleanerImpl struct {
	eventLifespan time.Duration
	repository    event.Repository
}

func NewCleanerImpl(eventLifespan time.Duration, repository event.Repository) *CleanerImpl {
	return &CleanerImpl{
		eventLifespan: eventLifespan,
		repository:    repository,
	}
}

func (c *CleanerImpl) Clean(ctx context.Context) error {
	maxEndDate := time.Now().Add(-c.eventLifespan)
	err := c.repository.DeleteWhereEndDateBefore(ctx, maxEndDate)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
