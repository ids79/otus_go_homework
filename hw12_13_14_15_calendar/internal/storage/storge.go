package storage

import (
	"context"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage/types"
	uuid "github.com/satori/go.uuid"
)

type Storage interface {
	Create(ev types.Event) (uuid.UUID, error)
	Update(u uuid.UUID, ev types.Event) error
	Delete(u uuid.UUID) error
	ListOnDay(time time.Time) []types.Event
	ListOnWeek(time time.Time) []types.Event
	ListOnMonth(time time.Time) []types.Event
}

func New(ctx context.Context, logg logger.Logg, config config.Config) Storage {
	switch config.Database.Storage {
	case "memory":
		return memorystorage.New()
	case "sql":
		stor := sqlstorage.New(ctx, logg, config)
		err := stor.Connect()
		if err != nil {
			return nil
		}
		err = stor.Migration()
		if err != nil {
			return nil
		}
		go func() {
			<-ctx.Done()
			stor.Close()
		}()
		return stor
	}
	logg.Error("Wrong type of storage")
	return nil
}
