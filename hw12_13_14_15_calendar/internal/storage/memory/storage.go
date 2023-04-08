package memory

import (
	"sync"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage/types"
	uuid "github.com/satori/go.uuid"
)

type Storage struct {
	sync.RWMutex
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

func (st *Storage) Create(ev types.Event) (uuid.UUID, error) {
	st.Lock()
	defer st.Unlock()
	ev.Year, ev.Month, ev.Day = ev.DateTime.Date()
	for _, e := range st.messages {
		if e.Year == ev.Year && e.Month == ev.Month && e.Day == ev.Day {
			return uuid.Nil, types.ErrDeteIsOccupied
		}
	}
	_, ev.Week = ev.DateTime.ISOWeek()
	u := uuid.NewV4()
	st.messages = append(st.messages, ev)
	st.mesID[u] = &ev
	return u, nil
}

func (st *Storage) GetEvent(u uuid.UUID) (*types.Event, error) {
	if ev, ok := st.mesID[u]; ok {
		return ev, nil
	}
	return &types.Event{}, types.ErrNotExistUUID
}

func (st *Storage) Update(u uuid.UUID, ev types.Event) error {
	st.Lock()
	defer st.Unlock()
	if _, ok := st.mesID[u]; !ok {
		return types.ErrNotExistUUID
	}
	st.mesID[u].Duration = ev.Duration
	st.mesID[u].Description = ev.Description
	st.mesID[u].TimeBefore = ev.TimeBefore
	return nil
}

func (st *Storage) Delete(u uuid.UUID) error {
	st.Lock()
	defer st.Unlock()
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

func (st *Storage) ListOnDay(time time.Time) []types.Event {
	list := make([]types.Event, 0)
	y, m, d := time.Date()
	for _, ev := range st.messages {
		if ev.Year == y && ev.Month == m && ev.Day == d {
			list = append(list, ev)
		}
	}
	return list
}

func (st *Storage) ListOnWeek(time time.Time) []types.Event {
	list := make([]types.Event, 0)
	y, w := time.ISOWeek()
	for _, ev := range st.messages {
		if ev.Year == y && ev.Week == w {
			list = append(list, ev)
		}
	}
	return list
}

func (st *Storage) ListOnMonth(time time.Time) []types.Event {
	list := make([]types.Event, 0)
	y, m, _ := time.Date()
	for _, ev := range st.messages {
		if ev.Year == y && ev.Month == m {
			list = append(list, ev)
		}
	}
	return list
}
