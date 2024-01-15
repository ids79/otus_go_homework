package memory

import (
	"context"
	"sync"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage/types"
	uuid "github.com/satori/go.uuid"
)

type Storage struct {
	sm       sync.RWMutex
	messages []types.Event
	mesID    map[uuid.UUID]*types.Event
}

func New() *Storage {
	messages := make([]types.Event, 0)
	mesID := make(map[uuid.UUID]*types.Event, 0)
	return &Storage{
		messages: messages,
		mesID:    mesID,
	}
}

func (st *Storage) Create(ctx context.Context, ev types.Event) (uuid.UUID, error) {
	st.sm.Lock()
	defer st.sm.Unlock()
	ev.Year, ev.Month, ev.Day = ev.DateTime.Date()
	for _, e := range st.messages {
		if e.Year == ev.Year && e.Month == ev.Month && e.Day == ev.Day {
			return uuid.Nil, types.ErrDateIsOccupied
		}
	}
	_, ev.Week = ev.DateTime.ISOWeek()
	u := uuid.NewV4()
	st.messages = append(st.messages, ev)
	st.mesID[u] = &ev
	return u, nil
}

func (st *Storage) GetEvent(u uuid.UUID) (*types.Event, error) {
	st.sm.Lock()
	defer st.sm.Unlock()
	if ev, ok := st.mesID[u]; ok {
		return ev, nil
	}
	return &types.Event{}, types.ErrNotExistUUID
}

func (st *Storage) Update(ctx context.Context, u uuid.UUID, ev types.Event) error {
	st.sm.Lock()
	defer st.sm.Unlock()
	if _, ok := st.mesID[u]; !ok {
		return types.ErrNotExistUUID
	}
	st.mesID[u].Duration = ev.Duration
	st.mesID[u].Description = ev.Description
	st.mesID[u].TimeBefore = ev.TimeBefore
	return nil
}

func (st *Storage) Delete(ctx context.Context, u uuid.UUID) error {
	st.sm.Lock()
	defer st.sm.Unlock()
	if _, ok := st.mesID[u]; !ok {
		return types.ErrNotExistUUID
	}
	for i, m := range st.messages {
		if m.ID == u {
			st.messages = append(st.messages[:i], st.messages[i+1:]...)
			break
		}
	}
	delete(st.mesID, u)
	return nil
}

func (st *Storage) ListOnDay(ctx context.Context, time time.Time) []types.Event {
	st.sm.Lock()
	defer st.sm.Unlock()
	list := make([]types.Event, 0)
	y, m, d := time.Date()
	for _, ev := range st.messages {
		if ev.Year == y && ev.Month == m && ev.Day == d {
			list = append(list, ev)
		}
	}
	return list
}

func (st *Storage) ListOnWeek(ctx context.Context, time time.Time) []types.Event {
	st.sm.Lock()
	defer st.sm.Unlock()
	list := make([]types.Event, 0)
	y, w := time.ISOWeek()
	for _, ev := range st.messages {
		if ev.Year == y && ev.Week == w {
			list = append(list, ev)
		}
	}
	return list
}

func (st *Storage) ListOnMonth(ctx context.Context, time time.Time) []types.Event {
	st.sm.Lock()
	defer st.sm.Unlock()
	list := make([]types.Event, 0)
	y, m, _ := time.Date()
	for _, ev := range st.messages {
		if ev.Year == y && ev.Month == m {
			list = append(list, ev)
		}
	}
	return list
}

func (st *Storage) SelectForReminder(ctx context.Context, t time.Time) []types.Event {
	return make([]types.Event, 0)
}

func (st *Storage) DeleteOldMessages(ctx context.Context, t time.Time) error {
	return nil
}
