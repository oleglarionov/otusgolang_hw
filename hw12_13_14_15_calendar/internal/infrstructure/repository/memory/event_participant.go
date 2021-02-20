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

func NewEventParticipantRepository() *EventParticipantRepository {
	return &EventParticipantRepository{
		participantsByEventID: make(map[event.ID][]event.Participant),
		participantsByUID:     make(map[user.UID][]event.Participant),
	}
}

func (r *EventParticipantRepository) Create(_ context.Context, participants []event.Participant) error {
	for _, participant := range participants {
		r.participantsByEventID[participant.EventID] = append(
			r.participantsByEventID[participant.EventID],
			participant,
		)

		r.participantsByUID[participant.UID] = append(
			r.participantsByUID[participant.UID],
			participant,
		)
	}

	return nil
}

func (r *EventParticipantRepository) GetParticipants(_ context.Context, eventID event.ID) ([]user.UID, error) {
	result := make([]user.UID, 0, len(r.participantsByEventID[eventID]))
	for _, p := range r.participantsByEventID[eventID] {
		result = append(result, p.UID)
	}

	return result, nil
}

func (r *EventParticipantRepository) GetUserEventIds(_ context.Context, uid user.UID) ([]event.ID, error) {
	result := make([]event.ID, 0)
	for _, p := range r.participantsByUID[uid] {
		result = append(result, p.EventID)
	}

	return result, nil
}

func (r *EventParticipantRepository) DeleteAllForEvent(ctx context.Context, eventID event.ID) error {
	participants := r.participantsByEventID[eventID]
	delete(r.participantsByEventID, eventID)

	for _, p := range participants {
		curSlice := r.participantsByUID[p.UID]
		for i, el := range curSlice {
			if el.EventID == eventID {
				r.participantsByUID[p.UID] = append(curSlice[:i], curSlice[i+1:]...)
				break
			}
		}
	}

	return nil
}
