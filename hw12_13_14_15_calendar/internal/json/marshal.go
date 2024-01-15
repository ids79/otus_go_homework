package internaljson

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	uuid "github.com/satori/go.uuid"
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	DateTime    JSONDate  `json:"datetime"`
	Duration    Duration  `json:"dur"`
	TimeBefore  Duration  `json:"timebefore"`
	Description string    `json:"desc"`
	UserID      string    `json:"user"`
}

type EventRem struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	DateTime    JSONDate  `json:"datetime"`
	Duration    Duration  `json:"dur"`
	Description string    `json:"desc"`
	UserID      string    `json:"user"`
}

type JSONDate time.Time

type Duration time.Duration

func (j JSONDate) MarshalJSON() ([]byte, error) {
	st := time.Time(j).Format("2006-01-02 03:04 PM")
	return json.Marshal(st)
}

func (j *JSONDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02 03:04 PM", s)
	if err != nil {
		return err
	}
	*j = JSONDate(t)
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	st := time.Duration(d).String()
	return json.Marshal(st)
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	t, err := time.ParseDuration(v)
	if err != nil {
		return err
	}
	*d = Duration(t)
	return nil
}

func EventsFormAppToView(eventsBase []app.Event) []Event {
	events := make([]Event, len(eventsBase))
	for i, ev := range eventsBase {
		events[i] = Event{
			ID:          ev.ID,
			Title:       ev.Title,
			DateTime:    JSONDate(ev.DateTime),
			Duration:    Duration(ev.Duration),
			Description: ev.Description,
			UserID:      strconv.Itoa(ev.UserID),
		}
	}
	return events
}

func EventsRemFormAppToView(eventsBase []app.Event) []EventRem {
	events := make([]EventRem, len(eventsBase))
	for i, ev := range eventsBase {
		events[i] = EventRem{
			ID:          ev.ID,
			Title:       ev.Title,
			DateTime:    JSONDate(ev.DateTime),
			Duration:    Duration(ev.Duration),
			Description: ev.Description,
			UserID:      strconv.Itoa(ev.UserID),
		}
	}
	return events
}
