package usecase

import (
	"context"
	"time"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
)

type EventDto struct {
	Title       string
	Description string
	BeginDate   time.Time
	EndDate     time.Time
}

type CreateEventDto struct {
	ID string
	EventDto
}

type UpdateEventDto struct {
	EventDto
}

type ReturnEventDto struct {
	event.Model
}

type EventUseCase interface {
	Create(ctx context.Context, uid user.UID, dto CreateEventDto) (*ReturnEventDto, error)
	Update(ctx context.Context, uid user.UID, id event.ID, dto UpdateEventDto) (*ReturnEventDto, error)
	Delete(ctx context.Context, uid user.UID, id event.ID) error
	DayList(ctx context.Context, uid user.UID, day time.Time) ([]*ReturnEventDto, error)
	WeekList(ctx context.Context, uid user.UID, beginDate time.Time) ([]*ReturnEventDto, error)
	MonthList(ctx context.Context, uid user.UID, beginDate time.Time) ([]*ReturnEventDto, error)
}

func NewEventUseCaseImpl(
	eventRepo event.Repository,
	participantRepo event.ParticipantRepository,
	service event.Service,
	uuidGenerator common.UUIDGenerator,
) *EventUseCaseImpl {
	return &EventUseCaseImpl{
		eventRepo:       eventRepo,
		participantRepo: participantRepo,
		service:         service,
		uuidGenerator:   uuidGenerator,
	}
}

var _ EventUseCase = (*EventUseCaseImpl)(nil)

type EventUseCaseImpl struct {
	eventRepo       event.Repository
	participantRepo event.ParticipantRepository
	service         event.Service
	uuidGenerator   common.UUIDGenerator
}

func (u *EventUseCaseImpl) Create(ctx context.Context, uid user.UID, dto CreateEventDto) (*ReturnEventDto, error) {
	// todo: добавить проверку, что beginDate < endDate

	err := u.service.EnsureIntervalAvailable(ctx, event.UserInterval{
		Interval: event.Interval{
			BeginDate: dto.BeginDate,
			EndDate:   dto.EndDate,
		},
		UID: uid,
	})
	if err != nil {
		return nil, err
	}

	id := dto.ID
	if id == "" {
		id = u.uuidGenerator.Generate()
	}

	model := event.Model{
		ID:          event.ID(id),
		Title:       dto.Title,
		Description: dto.Description,
		BeginDate:   dto.BeginDate,
		EndDate:     dto.EndDate,
	}

	err = u.service.CreateEvent(ctx, model, []user.UID{uid}) // todo: обернуть в транзакцию
	if err != nil {
		return nil, err
	}

	return &ReturnEventDto{model}, nil
}

func (u *EventUseCaseImpl) Update(ctx context.Context, uid user.UID, id event.ID, dto UpdateEventDto) (*ReturnEventDto, error) {
	if !u.service.HasAccess(ctx, uid, id) {
		return nil, ErrNotFound
	}

	model, err := u.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.service.EnsureIntervalAvailable(ctx, event.UserInterval{
		Interval: event.Interval{
			BeginDate: dto.BeginDate,
			EndDate:   dto.EndDate,
		},
		UID: uid,
	}, id)
	if err != nil {
		return nil, err
	}

	model.Title = dto.Title
	model.Description = dto.Description
	model.EndDate = dto.EndDate
	model.BeginDate = dto.BeginDate

	err = u.eventRepo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return &ReturnEventDto{model}, nil
}

func (u *EventUseCaseImpl) Delete(ctx context.Context, uid user.UID, id event.ID) error {
	if !u.service.HasAccess(ctx, uid, id) {
		return ErrNotFound
	}

	model, err := u.eventRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = u.service.DeleteEvent(ctx, model) // todo: обернуть в транзакцию
	if err != nil {
		return err
	}

	return nil
}

func (u *EventUseCaseImpl) DayList(ctx context.Context, uid user.UID, day time.Time) ([]*ReturnEventDto, error) {
	beginDate := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	endDate := beginDate.Add(time.Hour * 24)

	return u.intervalList(ctx, uid, beginDate, endDate)
}

func (u *EventUseCaseImpl) WeekList(ctx context.Context, uid user.UID, beginDate time.Time) ([]*ReturnEventDto, error) {
	beginDate = time.Date(beginDate.Year(), beginDate.Month(), beginDate.Day(), 0, 0, 0, 0, beginDate.Location())
	endDate := beginDate.AddDate(0, 0, 7)

	return u.intervalList(ctx, uid, beginDate, endDate)
}

func (u *EventUseCaseImpl) MonthList(ctx context.Context, uid user.UID, beginDate time.Time) ([]*ReturnEventDto, error) {
	beginDate = time.Date(beginDate.Year(), beginDate.Month(), beginDate.Day(), 0, 0, 0, 0, beginDate.Location())
	endDate := beginDate.AddDate(0, 1, 0)

	return u.intervalList(ctx, uid, beginDate, endDate)
}

func (u *EventUseCaseImpl) intervalList(
	ctx context.Context,
	uid user.UID,
	beginDate,
	endDate time.Time,
) ([]*ReturnEventDto, error) {
	models, err := u.eventRepo.GetByInterval(ctx, event.UserInterval{
		UID: uid,
		Interval: event.Interval{
			BeginDate: beginDate,
			EndDate:   endDate,
		},
	})
	if err != nil {
		return nil, err
	}

	result := make([]*ReturnEventDto, 0, len(models))
	for _, model := range models {
		result = append(result, &ReturnEventDto{model})
	}

	return result, nil
}
