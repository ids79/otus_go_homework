package storage

import (
	"context"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage/sql"
	typesevents "github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage/types-events"
	uuid "github.com/satori/go.uuid"
)

type Storage interface {
	Create(ctx context.Context, ev typesevents.Event) (uuid.UUID, error)
	Close() error
	Update(ctx context.Context, u uuid.UUID, ev typesevents.Event) error
	Delete(ctx context.Context, u uuid.UUID) error
	ListOnDay(ctx context.Context, time time.Time) []typesevents.Event
	ListOnWeek(ctx context.Context, time time.Time) []typesevents.Event
	ListOnMonth(ctx context.Context, time time.Time) []typesevents.Event
	SelectForReminder(ctx context.Context, time time.Time) []typesevents.Event
	DeleteOldMessages(ctx context.Context, t time.Time) error
}

func New(ctx context.Context, logg logger.Logg, config config.Config) Storage {
	switch config.Database.Storage {
	case "memory":
		return memorystorage.New()
	case "sql":
		stor := sqlstorage.New(logg, config)
		err := stor.Connect()
		if err != nil {
			return nil
		}
		err = stor.Migration()
		if err != nil {
			return nil
		}
		return stor
	}
	logg.Error("Wrong type of storage")
	return nil
}
