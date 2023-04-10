package sql

import (
	"context"
	"embed"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage/types"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	uuid "github.com/satori/go.uuid"
)

type Storage struct {
	ctx     context.Context
	connStr string
	logg    logger.Logg
	conn    *sqlx.DB
}

func New(ctx context.Context, logger logger.Logg, config config.Config) *Storage {
	return &Storage{
		ctx:     ctx,
		connStr: config.Database.ConnectString,
		logg:    logger,
	}
}

func (st *Storage) Connect() (err error) {
	st.conn, err = sqlx.ConnectContext(st.ctx, "pgx", st.connStr)
	if err != nil {
		st.logg.Error("cannot connect to base psql: ", err)
		return err
	}
	return st.conn.PingContext(st.ctx)
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (st *Storage) Migration() error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		st.logg.Error("Data migration failed with an error: ", err)
		return err
	}
	if err := goose.Up(st.conn.DB, "migrations"); err != nil {
		st.logg.Error("Data migration failed with an error: ", err)
		return err
	}
	st.logg.Info("Data migration was successful")
	return nil
}

func (st *Storage) Close() error {
	return st.conn.Close()
}

func (st *Storage) Create(ev types.Event) (uuid.UUID, error) {
	u := uuid.NewV4()
	y, m, d := ev.DateTime.Date()
	_, w := ev.DateTime.ISOWeek()
	query := `insert into events
	(id, title, description, date_time, year, month, week, day, duration, iser_id, time_before) 
	values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := st.conn.ExecContext(st.ctx, query,
		u.String(), ev.Title, ev.Description, ev.DateTime, y, int(m), w, d, int(ev.Duration), ev.UserID, int(ev.TimeBefore))
	if err != nil {
		return uuid.Nil, err
	}
	return u, nil
}

func (st *Storage) Update(u uuid.UUID, ev types.Event) error {
	query := `update events set description = $1, duration = $2, time_before = $3 where id = $4`
	_, err := st.conn.ExecContext(st.ctx, query,
		ev.Description, int(ev.Duration), int(ev.TimeBefore), u.String())
	if err != nil {
		return err
	}
	return nil
}

func (st *Storage) Delete(u uuid.UUID) error {
	query := `delete from events where id = $4`
	_, err := st.conn.ExecContext(st.ctx, query, u)
	if err != nil {
		return err
	}
	return nil
}

func getRows(logg logger.Logg, rows *sqlx.Rows) []types.Event {
	list := make([]types.Event, 0)
	for rows.Next() {
		var event types.Event
		err := rows.StructScan(&event)
		if err != nil {
			logg.Error("error in the request request processing", err)
			return nil
		}
		list = append(list, event)
	}
	return list
}

func (st *Storage) ListOnDay(time time.Time) []types.Event {
	y, m, d := time.Date()
	sql := "select * from events where day = :day and month = :month and year = :year"
	rows, err := st.conn.NamedQueryContext(st.ctx, sql, map[string]interface{}{
		"day":   d,
		"month": m,
		"year":  y,
	})
	if err != nil {
		st.logg.Error("error in the request formation process", err)
		return nil
	}
	return getRows(st.logg, rows)
}

func (st *Storage) ListOnWeek(time time.Time) []types.Event {
	y, w := time.ISOWeek()
	sql := "select * from events where week = :week and year = :year"
	rows, err := st.conn.NamedQueryContext(st.ctx, sql, map[string]interface{}{
		"week": w,
		"year": y,
	})
	if err != nil {
		st.logg.Error("error in the request formation process", err)
		return nil
	}
	return getRows(st.logg, rows)
}

func (st *Storage) ListOnMonth(time time.Time) []types.Event {
	y, m, _ := time.Date()
	sql := "select * from events where month = :month and year = :year"
	rows, err := st.conn.NamedQueryContext(st.ctx, sql, map[string]interface{}{
		"month": m,
		"year":  y,
	})
	if err != nil {
		st.logg.Error("error in the request formation process", err)
		return nil
	}
	return getRows(st.logg, rows)
}
