package types

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

var ErrNotExistUUID = errors.New("the event with this uuid is not found")

var ErrDeteIsOccupied = errors.New("the current date and time is occupied by another event")

type Event struct {
	ID          uuid.UUID
	Title       string
	DateTime    time.Time `db:"date_time"`
	Year        int
	Day         int
	Week        int
	Month       time.Month
	Duration    time.Duration
	Description string
	UserID      int           `db:"user_id"`
	TimeBefore  time.Duration `db:"time_before"`
}

type Notification struct {
	ID         uuid.UUID
	Title      string
	Date       time.Time
	UserDestID int
}
