package usecase

import (
	"context"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/repository/memory"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var uid = user.UID("user-1")
var now = time.Now()
var today = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
var yesterday = today.AddDate(0, 0, -1)
var tomorrow = today.AddDate(0, 0, 1)
var weekLater = today.AddDate(0, 0, 7)
var monthLater = today.AddDate(0, 1, 0)

func TestEventUseCaseImpl_Create(t *testing.T) {
	t.Run("positive create", func(t *testing.T) {
		uc := buildEventUseCase(initData{})
		eventID := "fcbc5856-528a-11eb-a567-acde48001122"
		beginDate := now
		endDate := now.Add(time.Hour)

		result, err := uc.Create(context.Background(), uid, CreateEventDto{
			ID: eventID,
			EventDto: EventDto{
				Title:       "title-1",
				Description: "descr-1",
				BeginDate:   beginDate,
				EndDate:     endDate,
			},
		})

		expectedModel := event.Model{
			ID:          event.ID(eventID),
			Title:       "title-1",
			Description: "descr-1",
			BeginDate:   beginDate,
			EndDate:     endDate,
		}

		require.NoError(t, err)
		require.Equal(t, &ReturnEventDto{
			Model: expectedModel,
		}, result)

		model, err := uc.eventRepo.GetByID(context.Background(), event.ID(eventID))
		require.NoError(t, err)
		require.Equal(t, expectedModel, model)

		participants, err := uc.participantRepo.GetParticipants(context.Background(), event.ID(eventID))
		require.NoError(t, err)
		require.Equal(t, []user.UID{uid}, participants)
	})
	t.Run("create with unavailable interval", func(t *testing.T) {
		uc := buildEventUseCase(initData{
			events: []event.Model{{
				ID:          "id1",
				Title:       "t1",
				Description: "d1",
				BeginDate:   now,
				EndDate:     now.Add(time.Hour),
			}},
			participants: []event.Participant{{
				EventID: "id1",
				UID:     uid,
			}},
		})

		_, err := uc.Create(context.Background(), uid, CreateEventDto{
			ID: "id2",
			EventDto: EventDto{
				Title:       "t2",
				Description: "d2",
				BeginDate:   now.Add(30 * time.Minute),
				EndDate:     now.Add(2 * time.Hour),
			},
		})
		require.Equal(t, event.ErrIntervalNotAvailable, err)
	})
}

func TestEventUseCaseImpl_Update(t *testing.T) {
	t.Run("positive update", func(t *testing.T) {
		eventID := event.ID("fcbc5856-528a-11eb-a567-acde48001122")
		uc := buildEventUseCase(initData{
			events: []event.Model{
				{
					ID:          eventID,
					Title:       "title-1",
					Description: "descr-1",
					BeginDate:   now,
					EndDate:     now.Add(time.Hour),
				},
			},
			participants: []event.Participant{
				{
					EventID: eventID,
					UID:     uid,
				},
			},
		})

		result, err := uc.Update(context.Background(), uid, eventID, UpdateEventDto{
			EventDto: EventDto{
				Title:       "title-1-updated",
				Description: "descr-1-updated",
				BeginDate:   now.Add(time.Hour),
				EndDate:     now.Add(2 * time.Hour),
			},
		})

		expectedModel := event.Model{
			ID:          eventID,
			Title:       "title-1-updated",
			Description: "descr-1-updated",
			BeginDate:   now.Add(time.Hour),
			EndDate:     now.Add(2 * time.Hour),
		}

		require.NoError(t, err)
		require.Equal(t, &ReturnEventDto{expectedModel}, result)

		model, err := uc.eventRepo.GetByID(context.Background(), eventID)
		require.NoError(t, err)
		require.Equal(t, expectedModel, model)

		participants, err := uc.participantRepo.GetParticipants(context.Background(), eventID)
		require.NoError(t, err)
		require.Equal(t, []user.UID{uid}, participants)
	})
	t.Run("update to unavailable interval", func(t *testing.T) {
		uc := buildEventUseCase(initData{
			events: []event.Model{
				{
					ID:          "id1",
					Title:       "t1",
					Description: "d1",
					BeginDate:   now,
					EndDate:     now.Add(time.Hour),
				}, {
					ID:          "id2",
					Title:       "t2",
					Description: "d2",
					BeginDate:   now.Add(time.Hour),
					EndDate:     now.Add(2 * time.Hour),
				},
			},
			participants: []event.Participant{
				{
					EventID: "id1",
					UID:     uid,
				}, {
					EventID: "id2",
					UID:     uid,
				},
			},
		})

		_, err := uc.Update(context.Background(), uid, "id2", UpdateEventDto{
			EventDto: EventDto{
				Title:       "t21",
				Description: "d21",
				BeginDate:   now.Add(30 * time.Minute),
				EndDate:     now.Add(2 * time.Hour),
			},
		})

		require.Equal(t, event.ErrIntervalNotAvailable, err)
	})
}

func TestEventUseCaseImpl_Delete(t *testing.T) {
	now := time.Now()

	t.Run("positive delete", func(t *testing.T) {
		uid := user.UID("user-1")
		eventID := event.ID("fcbc5856-528a-11eb-a567-acde48001122")
		uc := buildEventUseCase(initData{
			events: []event.Model{
				{
					ID:          eventID,
					Title:       "title-1",
					Description: "descr-1",
					BeginDate:   now,
					EndDate:     now.Add(time.Hour),
				},
			},
			participants: []event.Participant{
				{
					EventID: eventID,
					UID:     uid,
				},
			},
		})

		err := uc.Delete(context.Background(), uid, eventID)
		require.NoError(t, err)

		model, err := uc.eventRepo.GetByID(context.Background(), eventID)
		require.Error(t, err)
		require.Equal(t, event.Model{}, model)

		participants, err := uc.participantRepo.GetParticipants(context.Background(), eventID)
		require.NoError(t, err)
		require.Equal(t, []user.UID{}, participants)

	})
}

func TestEventUseCaseImpl_DayList(t *testing.T) {
	t.Run("day list", func(t *testing.T) {
		eventID1 := event.ID("fcbc5856-528a-11eb-a567-acde48001122")
		eventID2 := event.ID("fcbc5856-528a-11eb-a567-acde48001123")
		eventID3 := event.ID("fcbc5856-528a-11eb-a567-acde48001124")
		uid := user.UID("user-1")

		uc := buildEventUseCase(initData{
			events: []event.Model{
				{
					ID:          eventID1,
					Title:       "title-1",
					Description: "descr-1",
					BeginDate:   today,
					EndDate:     today.Add(time.Hour),
				},
				{
					ID:          eventID2,
					Title:       "title-2",
					Description: "descr-2",
					BeginDate:   tomorrow,
					EndDate:     tomorrow.Add(25 * time.Hour),
				},
				{
					ID:          eventID3,
					Title:       "title-3",
					Description: "descr-3",
					BeginDate:   yesterday,
					EndDate:     yesterday.Add(time.Hour),
				},
			},
			participants: []event.Participant{
				{
					EventID: eventID1,
					UID:     uid,
				},
				{
					EventID: eventID2,
					UID:     uid,
				},
			},
		})

		result, err := uc.DayList(context.Background(), uid, today)
		require.NoError(t, err)
		require.Equal(t, []*ReturnEventDto{
			{
				Model: event.Model{
					ID:          eventID1,
					Title:       "title-1",
					Description: "descr-1",
					BeginDate:   today,
					EndDate:     today.Add(time.Hour),
				},
			},
		}, result)
	})
}

func TestEventUseCaseImpl_WeekList(t *testing.T) {
	t.Run("month list", func(t *testing.T) {
		eventID1 := event.ID("fcbc5856-528a-11eb-a567-acde48001122")
		eventID2 := event.ID("fcbc5856-528a-11eb-a567-acde48001123")
		uid := user.UID("user-1")

		uc := buildEventUseCase(initData{
			events: []event.Model{
				{
					ID:          eventID1,
					Title:       "title-1",
					Description: "descr-1",
					BeginDate:   today,
					EndDate:     today.Add(time.Hour),
				},
				{
					ID:          eventID2,
					Title:       "title-2",
					Description: "descr-2",
					BeginDate:   weekLater,
					EndDate:     weekLater.Add(time.Hour),
				},
			},
			participants: []event.Participant{
				{
					EventID: eventID1,
					UID:     uid,
				},
				{
					EventID: eventID2,
					UID:     uid,
				},
			},
		})

		result, err := uc.DayList(context.Background(), uid, today)
		require.NoError(t, err)
		require.Equal(t, []*ReturnEventDto{
			{
				Model: event.Model{
					ID:          eventID1,
					Title:       "title-1",
					Description: "descr-1",
					BeginDate:   today,
					EndDate:     today.Add(time.Hour),
				},
			},
		}, result)
	})
}

func TestEventUseCaseImpl_MonthList(t *testing.T) {
	t.Run("month list", func(t *testing.T) {
		eventID1 := event.ID("fcbc5856-528a-11eb-a567-acde48001122")
		eventID2 := event.ID("fcbc5856-528a-11eb-a567-acde48001123")
		uid := user.UID("user-1")

		uc := buildEventUseCase(initData{
			events: []event.Model{
				{
					ID:          eventID1,
					Title:       "title-1",
					Description: "descr-1",
					BeginDate:   today,
					EndDate:     today.Add(time.Hour),
				},
				{
					ID:          eventID2,
					Title:       "title-2",
					Description: "descr-2",
					BeginDate:   monthLater,
					EndDate:     monthLater.Add(time.Hour),
				},
			},
			participants: []event.Participant{
				{
					EventID: eventID1,
					UID:     uid,
				},
				{
					EventID: eventID2,
					UID:     uid,
				},
			},
		})

		result, err := uc.MonthList(context.Background(), uid, today)
		require.NoError(t, err)
		require.Equal(t, []*ReturnEventDto{
			{
				Model: event.Model{
					ID:          eventID1,
					Title:       "title-1",
					Description: "descr-1",
					BeginDate:   today,
					EndDate:     today.Add(time.Hour),
				},
			},
		}, result)
	})
}

type initData struct {
	events       []event.Model
	participants []event.Participant
}

func buildEventUseCase(data initData) *EventUseCaseImpl {
	ctx := context.Background()

	participantRepository := memory.NewEventParticipantRepository()
	eventRepository := memory.NewEventRepository(participantRepository)
	eventService := event.NewService(eventRepository, participantRepository)
	uuidGenerator := uuid.NewGenerator()

	for _, e := range data.events {
		_ = eventRepository.Create(ctx, e)
	}
	_ = participantRepository.Create(ctx, data.participants)

	return NewEventUseCaseImpl(eventRepository, participantRepository, eventService, uuidGenerator)
}
