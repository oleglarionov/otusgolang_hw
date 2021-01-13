package handler

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/api"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventServiceServerImpl struct {
	api.UnimplementedEventServiceServer
	eventUseCase usecase.EventUseCase
}

func NewEventServiceServerImpl(eventUseCase usecase.EventUseCase) *EventServiceServerImpl {
	return &EventServiceServerImpl{
		eventUseCase: eventUseCase,
	}
}

func (s *EventServiceServerImpl) Create(ctx context.Context, request *api.CreateEventRequest) (*api.Event, error) {
	result, err := s.eventUseCase.Create(ctx, getUID(ctx), usecase.CreateEventDto{
		ID: request.Id,
		EventDto: usecase.EventDto{
			Title:       request.Title,
			Description: request.Description,
			BeginDate:   request.BeginDate.AsTime(),
			EndDate:     request.EndDate.AsTime(),
		},
	})
	if err != nil {
		return nil, err
	}

	return toAPIEvent(result), nil
}

func (s *EventServiceServerImpl) Update(ctx context.Context, request *api.UpdateEventRequest) (*api.Event, error) {
	result, err := s.eventUseCase.Update(ctx, getUID(ctx), event.ID(request.Id), usecase.UpdateEventDto{
		EventDto: usecase.EventDto{
			Title:       request.Title,
			Description: request.Description,
			BeginDate:   request.BeginDate.AsTime(),
			EndDate:     request.EndDate.AsTime(),
		},
	})
	if err != nil {
		return nil, err
	}

	return toAPIEvent(result), nil
}

func (s *EventServiceServerImpl) Delete(ctx context.Context, request *api.DeleteEventRequest) (*empty.Empty, error) {
	err := s.eventUseCase.Delete(ctx, getUID(ctx), event.ID(request.Id))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *EventServiceServerImpl) DayList(ctx context.Context, request *api.DayListRequest) (*api.Events, error) {
	result, err := s.eventUseCase.DayList(ctx, getUID(ctx), request.Day.AsTime())
	if err != nil {
		return nil, err
	}

	return toAPIEvents(result), nil
}

func (s *EventServiceServerImpl) WeekList(ctx context.Context, request *api.WeekListRequest) (*api.Events, error) {
	result, err := s.eventUseCase.WeekList(ctx, getUID(ctx), request.BeginDate.AsTime())
	if err != nil {
		return nil, err
	}

	return toAPIEvents(result), nil
}

func (s *EventServiceServerImpl) MonthList(ctx context.Context, request *api.MonthListRequest) (*api.Events, error) {
	result, err := s.eventUseCase.MonthList(ctx, getUID(ctx), request.BeginDate.AsTime())
	if err != nil {
		return nil, err
	}

	return toAPIEvents(result), nil
}

func getUID(ctx context.Context) user.UID {
	return ctx.Value(UIDKey{}).(user.UID)
}

func toAPIEvent(dto *usecase.ReturnEventDto) *api.Event {
	return &api.Event{
		Id:          string(dto.ID),
		Title:       dto.Title,
		Description: dto.Description,
		BeginDate:   timestamppb.New(dto.BeginDate),
		EndDate:     timestamppb.New(dto.EndDate),
	}
}

func toAPIEvents(dtoList []*usecase.ReturnEventDto) *api.Events {
	result := make([]*api.Event, 0, len(dtoList))
	for _, dto := range dtoList {
		result = append(result, toAPIEvent(dto))
	}

	return &api.Events{Items: result}
}
