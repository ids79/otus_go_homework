package app

import (
	"context"
	"fmt"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage/types"
	uuid "github.com/satori/go.uuid"
)

type App struct {
	logg    logger.Logg
	storage storage.Storage
	conf    config.Config
}

type Application interface {
	CreateEvent(ctx context.Context, ev Event) uuid.UUID
	UpgateEvent(ctx context.Context, u uuid.UUID, ev Event) error
	DeleteEvent(ctx context.Context, u uuid.UUID) error
	GetListOnDay(ctx context.Context, time time.Time) []Event
	GetListOnWeek(ctx context.Context, time time.Time) []Event
	GetListOnMonth(ctx context.Context, time time.Time) []Event
	SelectForReminder(ctx context.Context, time time.Time) []Event
	DeleteOldMessages(ctx context.Context, t time.Time) error
}

type Event struct {
	ID          uuid.UUID
	Title       string
	DateTime    time.Time
	Duration    time.Duration
	TimeBefore  time.Duration
	Description string
	UserID      int
}

func New(logger logger.Logg, storage storage.Storage, config config.Config) Application {
	return &App{
		logg:    logger,
		storage: storage,
		conf:    config,
	}
}

func (a *App) CreateEvent(ctx context.Context, ev Event) uuid.UUID {
	u, err := a.storage.Create(ctx, types.Event{
		DateTime:    ev.DateTime,
		Title:       ev.Title,
		Duration:    ev.Duration,
		Description: ev.Description,
		TimeBefore:  ev.TimeBefore,
		UserID:      ev.UserID,
	})
	if err != nil {
		a.logg.Error("error with adding an event: ", err)
		return uuid.Nil
	}
	return u
}

func (a *App) UpgateEvent(ctx context.Context, u uuid.UUID, ev Event) error {
	err := a.storage.Update(ctx, u, types.Event{
		Duration:    ev.Duration,
		Description: ev.Description,
		TimeBefore:  ev.TimeBefore,
	})
	if err != nil {
		a.logg.Error("error with update event: ", u, " ", err)
	}
	return err
}

func (a *App) DeleteEvent(ctx context.Context, u uuid.UUID) error {
	err := a.storage.Delete(ctx, u)
	if err != nil {
		a.logg.Error("error with delete event: ", u, " ", err)
	}
	return err
}

func eventsFormBaseToApp(eventsBase []types.Event) []Event {
	events := make([]Event, len(eventsBase))
	for i, ev := range eventsBase {
		events[i] = Event{
			ID:          ev.ID,
			Title:       ev.Title,
			DateTime:    ev.DateTime,
			Duration:    ev.Duration,
			Description: ev.Description,
			UserID:      ev.UserID,
		}
	}
	return events
}

func (a *App) GetListOnDay(ctx context.Context, time time.Time) []Event {
	events := a.storage.ListOnDay(ctx, time)
	if events == nil {
		return nil
	}
	if len(events) == 0 {
		a.logg.Info(fmt.Sprintf("nothing was selected on the day %s", time.Format("2006-01-02")))
		return make([]Event, 0)
	}
	return eventsFormBaseToApp(events)
}

func (a *App) GetListOnWeek(ctx context.Context, time time.Time) []Event {
	events := a.storage.ListOnWeek(ctx, time)
	if events == nil {
		return nil
	}
	if len(events) == 0 {
		_, w := time.ISOWeek()
		a.logg.Info(fmt.Sprintf("nothing was selected on the week %d", w))
		return make([]Event, 0)
	}
	return eventsFormBaseToApp(events)
}

func (a *App) GetListOnMonth(ctx context.Context, time time.Time) []Event {
	events := a.storage.ListOnMonth(ctx, time)
	if events == nil {
		return nil
	}
	if len(events) == 0 {
		a.logg.Info(fmt.Sprintf("nothing was selected on the month %s", time.Format("2006-01")))
		return make([]Event, 0)
	}
	return eventsFormBaseToApp(events)
}

func (a *App) SelectForReminder(ctx context.Context, time time.Time) []Event {
	events := a.storage.SelectForReminder(ctx, time)
	if events == nil {
		return nil
	}
	if len(events) == 0 {
		return make([]Event, 0)
	}
	return eventsFormBaseToApp(events)
}

func (a *App) DeleteOldMessages(ctx context.Context, t time.Time) error {
	return a.storage.DeleteOldMessages(ctx, t)
}
