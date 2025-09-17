package memory

import (
	"context"
	"testing"
	"time"

	typesevents "github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage/types-events"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	t.Run("base test", func(t *testing.T) {
		storage := New()

		u, err := storage.Create(ctx, typesevents.Event{
			DateTime:    time.Now(),
			Title:       "Event",
			Duration:    time.Hour,
			Description: "New event ...",
			TimeBefore:  time.Hour * 6,
			UserID:      1,
		})
		require.NoError(t, err)
		require.NotNil(t, u)
		ev, err := storage.GetEvent(u)
		require.NoError(t, err)
		require.Equal(t, "Event", ev.Title)
		require.Equal(t, "New event ...", ev.Description)
		require.Equal(t, time.Hour, ev.Duration)
		require.Equal(t, time.Hour*6, ev.TimeBefore)
		require.Equal(t, 1, ev.UserID)

		_, err = storage.Create(ctx, typesevents.Event{
			DateTime:    time.Now(),
			Title:       "Event",
			Duration:    time.Hour,
			Description: "New event ...",
			TimeBefore:  time.Hour * 6,
			UserID:      1,
		})
		require.ErrorIs(t, typesevents.ErrDateIsOccupied, err)

		err = storage.Update(ctx, u, typesevents.Event{
			Duration:    time.Hour * 2,
			Description: "Change event ...",
			TimeBefore:  time.Hour * 12,
		})
		require.NoError(t, err)
		require.Equal(t, "Change event ...", ev.Description)
		require.Equal(t, time.Hour*2, ev.Duration)
		require.Equal(t, time.Hour*12, ev.TimeBefore)

		err = storage.Update(ctx, uuid.NewV4(), typesevents.Event{})
		require.ErrorIs(t, typesevents.ErrNotExistUUID, err)

		events := storage.ListOnDay(ctx, time.Now())
		require.Equal(t, 1, len(events))

		err = storage.Delete(ctx, u)
		require.NoError(t, err)
		_, err = storage.GetEvent(u)
		require.ErrorIs(t, typesevents.ErrNotExistUUID, err)

		err = storage.Delete(ctx, u)
		require.ErrorIs(t, typesevents.ErrNotExistUUID, err)
	})
}
