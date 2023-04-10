package types

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

var ErrNotExistUUID = errors.New("the event with this uuid is not exist")

var ErrDeteIsOccupied = errors.New("the current date is occupied by another event")

type Event struct {
	ID          uuid.UUID
	Title       string
	DateTime    time.Time
	Year        int
	Day         int
	Week        int
	Month       time.Month
	Duration    time.Duration
	Description string
	UserID      int
	TimeBefore  time.Duration
}

type Notification struct {
	ID         uuid.UUID
	Title      string
	Date       time.Time
	UserDestID int
}
