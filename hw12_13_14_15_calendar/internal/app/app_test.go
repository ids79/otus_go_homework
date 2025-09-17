package app

import (
	"context"
	"testing"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	typesevents "github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage/types-events"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type StorageMock struct {
	mock.Mock
}

func (am *StorageMock) Create(ctx context.Context, ev typesevents.Event) (uuid.UUID, error) {
	_ = ctx
	t, _ := time.Parse("2006-01-02 03:04 PM", "2023-11-15 10:30 AM")
	if time.Time.Equal(ev.DateTime, t) {
		return uuid.Nil, typesevents.ErrDateIsOccupied
	}
	return uuid.NewV4(), nil
}

func (am *StorageMock) Update(ctx context.Context, u uuid.UUID, ev typesevents.Event) error {
	_ = ctx
	_ = ev
	uid, _ := uuid.FromString("02129661-9c49-48de-8d06-f88fe3867279")
	if u == uid {
		return typesevents.ErrNotExistUUID
	}
	return nil
}

func (am *StorageMock) Delete(ctx context.Context, u uuid.UUID) error {
	args := am.Called(ctx, u)
	return args.Error(0)
}

func (am *StorageMock) ListOnDay(ctx context.Context, tm time.Time) []typesevents.Event {
	_ = ctx
	y, m, d := tm.Date()
	if y != 2023 || m != 11 || d != 15 {
		return make([]typesevents.Event, 0)
	}
	ev := make([]typesevents.Event, 0)
	ev = append(ev, typesevents.Event{
		ID:          uuid.FromStringOrNil("02129661-9c49-48de-8d06-f88fe3867278"),
		DateTime:    tm,
		Title:       "Message",
		Duration:    600,
		Description: "Desc",
		TimeBefore:  7200,
		UserID:      1,
	})
	return ev
}

func (am *StorageMock) ListOnWeek(ctx context.Context, time time.Time) []typesevents.Event {
	_ = ctx
	_ = time
	return make([]typesevents.Event, 0)
}

func (am *StorageMock) ListOnMonth(ctx context.Context, time time.Time) []typesevents.Event {
	_ = ctx
	_ = time
	return make([]typesevents.Event, 0)
}

func (am *StorageMock) SelectForReminder(ctx context.Context, time time.Time) []typesevents.Event {
	_ = ctx
	_ = time
	return make([]typesevents.Event, 0)
}

func (am *StorageMock) DeleteOldMessages(ctx context.Context, t time.Time) error {
	_ = ctx
	_ = t
	return nil
}

func (am *StorageMock) Close() error {
	return nil
}

func TestStorage(t *testing.T) {
	ctx := context.Background()
	storageMock := &StorageMock{}
	storageMock.On("Delete", ctx, uuid.FromStringOrNil("02129661-9c49-48de-8d06-f88fe3867279")).
		Once().
		Return(typesevents.ErrNotExistUUID)
	storageMock.On("Delete", ctx, uuid.FromStringOrNil("02129661-9c49-48de-8d06-f88fe3867278")).Return(nil)

	tm, _ := time.Parse("2006-01-02", "2023-11-10")
	tm2, _ := time.Parse("2006-01-02 03:04 PM", "2023-11-15 10:30 AM")
	evTarget := make([]Event, 1)
	ev := Event{
		ID:          uuid.FromStringOrNil("02129661-9c49-48de-8d06-f88fe3867278"),
		DateTime:    tm,
		Title:       "Message",
		Duration:    600,
		Description: "Desc",
		TimeBefore:  7200,
		UserID:      1,
	}
	evTarget[0] = ev

	config := config.NewConfig("../../configs/calendar_config.toml")
	logg := logger.New(config.Logger, "Tests:")
	calendar := New(logg, storageMock, config)

	t.Run("base test", func(t *testing.T) {
		u := calendar.CreateEvent(ctx, ev)
		require.NotEqualValues(t, uuid.Nil, u)
		ev.DateTime = tm2
		u = calendar.CreateEvent(ctx, ev)
		require.EqualValues(t, uuid.Nil, u)

		err := calendar.UpgateEvent(ctx, uuid.FromStringOrNil("02129661-9c49-48de-8d06-f88fe3867279"), ev)
		require.ErrorIs(t, err, typesevents.ErrNotExistUUID)
		err = calendar.UpgateEvent(ctx, uuid.FromStringOrNil("02129661-9c49-48de-8d06-f88fe3867278"), ev)
		require.Nil(t, err)

		err = calendar.DeleteEvent(ctx, uuid.FromStringOrNil("02129661-9c49-48de-8d06-f88fe3867279"))
		require.ErrorIs(t, err, typesevents.ErrNotExistUUID)
		err = calendar.DeleteEvent(ctx, uuid.FromStringOrNil("02129661-9c49-48de-8d06-f88fe3867278"))
		require.Nil(t, err)

		evTarget[0].TimeBefore = 0
		evActual := calendar.GetListOnDay(ctx, tm)
		require.Empty(t, evActual)
		evTarget[0].DateTime = tm2
		evActual = calendar.GetListOnDay(ctx, tm2)
		require.EqualValues(t, evTarget, evActual)
	})
	storageMock.AssertExpectations(t)
}
