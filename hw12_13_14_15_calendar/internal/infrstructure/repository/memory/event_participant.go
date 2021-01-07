package memory

import (
	"context"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
)

type EventParticipantRepository struct {
	participantsByEventID map[event.ID][]event.Participant
	participantsByUID     map[user.UID][]event.Participant
}

func NewEventParticipantRepository() event.ParticipantRepository {
	return &EventParticipantRepository{
		participantsByEventID: make(map[event.ID][]event.Participant),
		participantsByUID:     make(map[user.UID][]event.Participant),
	}
}

func (r *EventParticipantRepository) Create(_ context.Context, participants []event.Participant) error {
	for _, participant := range participants {
		r.participantsByEventID[participant.EventId] = append(
			r.participantsByEventID[participant.EventId],
			participants...,
		)

		r.participantsByUID[participant.Uid] = append(
			r.participantsByUID[participant.Uid],
			participants...,
		)
	}

	return nil
}

func (r *EventParticipantRepository) GetParticipants(_ context.Context, eventId event.ID) ([]user.UID, error) {
	result := make([]user.UID, 0)
	for _, p := range r.participantsByEventID[eventId] {
		result = append(result, p.Uid)
	}

	return result, nil
}

func (r *EventParticipantRepository) GetUserEventIds(_ context.Context, uid user.UID) ([]event.ID, error) {
	result := make([]event.ID, 0)
	for _, p := range r.participantsByUID[uid] {
		result = append(result, p.EventId)
	}

	return result, nil
}

func (r *EventParticipantRepository) DeleteAllForEvent(ctx context.Context, eventId event.ID) error {
	participants := r.participantsByEventID[eventId]
	delete(r.participantsByEventID, eventId)

	for _, p := range participants {
		curSlice := r.participantsByUID[p.Uid]
		for i, el := range curSlice {
			if el.EventId == eventId {
				r.participantsByUID[p.Uid] = append(curSlice[:i], curSlice[i+1:]...)
				break
			}
		}
	}

	return nil
}
