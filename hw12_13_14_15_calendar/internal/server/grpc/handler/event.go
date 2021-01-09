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
	eventUseCase usecase.EventUseCaseInterface
}

func NewEventServiceServerImpl(eventUseCase usecase.EventUseCaseInterface) *EventServiceServerImpl {
	return &EventServiceServerImpl{
		eventUseCase: eventUseCase,
	}
}

func (s *EventServiceServerImpl) Create(ctx context.Context, request *api.CreateEventRequest) (*api.Event, error) {
	result, err := s.eventUseCase.Create(ctx, getUid(ctx), usecase.CreateEventDto{
		Id: request.Id,
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

	return toApiEvent(result), nil
}

func (s *EventServiceServerImpl) Update(ctx context.Context, request *api.UpdateEventRequest) (*api.Event, error) {
	result, err := s.eventUseCase.Update(ctx, getUid(ctx), event.ID(request.Id), usecase.UpdateEventDto{
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

	return toApiEvent(result), nil
}

func (s *EventServiceServerImpl) Delete(ctx context.Context, request *api.DeleteEventRequest) (*empty.Empty, error) {
	err := s.eventUseCase.Delete(ctx, getUid(ctx), event.ID(request.Id))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *EventServiceServerImpl) DayList(ctx context.Context, request *api.DayListRequest) (*api.Events, error) {
	result, err := s.eventUseCase.DayList(ctx, getUid(ctx), request.Day.AsTime())
	if err != nil {
		return nil, err
	}

	return toApiEvents(result), nil
}

func (s *EventServiceServerImpl) WeekList(ctx context.Context, request *api.WeekListRequest) (*api.Events, error) {
	result, err := s.eventUseCase.WeekList(ctx, getUid(ctx), request.BeginDate.AsTime())
	if err != nil {
		return nil, err
	}

	return toApiEvents(result), nil
}

func (s *EventServiceServerImpl) MonthList(ctx context.Context, request *api.MonthListRequest) (*api.Events, error) {
	result, err := s.eventUseCase.MonthList(ctx, getUid(ctx), request.BeginDate.AsTime())
	if err != nil {
		return nil, err
	}

	return toApiEvents(result), nil
}

func getUid(ctx context.Context) user.UID {
	return ctx.Value("uid").(user.UID)
}

func toApiEvent(dto *usecase.ReturnEventDto) *api.Event {
	return &api.Event{
		Id:          string(dto.Id),
		Title:       dto.Title,
		Description: dto.Description,
		BeginDate:   timestamppb.New(dto.BeginDate),
		EndDate:     timestamppb.New(dto.EndDate),
	}
}

func toApiEvents(dtoList []*usecase.ReturnEventDto) *api.Events {
	result := make([]*api.Event, 0, len(dtoList))
	for _, dto := range dtoList {
		result = append(result, toApiEvent(dto))
	}

	return &api.Events{Items: result}
}
