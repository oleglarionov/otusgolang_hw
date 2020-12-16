package memory

import (
	"context"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/require"
	"strconv"
	"sync"
	"testing"
)

func TestEventRepo(t *testing.T) {
	ctx := context.Background()
	initialEvents := []model.Event{
		{
			ID:    "1",
			Title: "event-1",
		},
		{
			ID:    "2",
			Title: "event-2",
		},
	}

	t.Run("create", func(t *testing.T) {
		repo := makeRepo(initialEvents)

		event := model.Event{
			ID:    "3",
			Title: "event-3",
		}

		err := repo.Create(ctx, event)

		require.NoError(t, err)
		require.Equal(t, event, repo.data[event.ID])
	})

	t.Run("get all", func(t *testing.T) {
		repo := makeRepo(initialEvents)

		events, err := repo.GetAll(ctx)
		require.NoError(t, err)

		require.Equal(t, initialEvents, events)
	})

	t.Run("get by id", func(t *testing.T) {
		repo := makeRepo(initialEvents)

		event, err := repo.GetByID(ctx, initialEvents[0].ID)
		require.NoError(t, err)

		require.Equal(t, initialEvents[0], event)
	})

	t.Run("update", func(t *testing.T) {
		repo := makeRepo(initialEvents)
		updatedEvent := model.Event{
			ID:    initialEvents[0].ID,
			Title: "updated-title",
		}

		err := repo.Update(ctx, updatedEvent)
		require.NoError(t, err)

		require.Equal(t, updatedEvent, repo.data[updatedEvent.ID])
		require.Equal(t, initialEvents[1], repo.data[initialEvents[1].ID])
	})

	t.Run("delete", func(t *testing.T) {
		repo := makeRepo(initialEvents)
		eventToDelete := initialEvents[0]

		err := repo.Delete(ctx, eventToDelete)
		require.NoError(t, err)

		_, ok := repo.data[eventToDelete.ID]
		require.False(t, ok)

		_, ok = repo.data[initialEvents[1].ID]
		require.True(t, ok)
	})

	t.Run("race", func(t *testing.T) {
		repo := makeRepo(nil)
		n := 20

		wg := sync.WaitGroup{}
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func(idInt int) {
				defer wg.Done()

				id := strconv.Itoa(idInt)
				err := repo.Create(ctx, model.Event{
					ID:    model.EventID(id),
					Title: "title-" + id,
				})

				require.NoError(t, err)
			}(i)
		}
		wg.Wait()

		require.Equal(t, n, len(repo.data))

		wg.Add(n)
		for i := 0; i < n; i++ {
			go func(idInt int) {
				defer wg.Done()

				id := model.EventID(strconv.Itoa(idInt))

				event, err := repo.GetByID(ctx, id)
				require.NoError(t, err)

				err = repo.Delete(ctx, event)
				require.NoError(t, err)
			}(i)
		}
		wg.Wait()

		require.Equal(t, 0, len(repo.data))
	})
}

func makeRepo(data []model.Event) *EventRepo {
	ctx := context.Background()

	repo := NewEventRepo()
	for _, event := range data {
		_ = repo.Create(ctx, event)
	}

	return repo
}
