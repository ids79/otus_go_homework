package sql

import (
	"context"
	"embed"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage/types"
	"github.com/jmoiron/sqlx"
	goose "github.com/pressly/goose/v3"
	uuid "github.com/satori/go.uuid"
)

type Storage struct {
	connStr string
	logg    logger.Logg
	conn    *sqlx.DB
}

func New(logger logger.Logg, config config.Config) *Storage {
	return &Storage{
		connStr: config.Database.ConnectString,
		logg:    logger,
	}
}

func (st *Storage) Connect() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	st.conn, err = sqlx.ConnectContext(ctx, "pgx", st.connStr)
	if err != nil {
		st.logg.Error("cannot connect to base psql: ", err)
		return err
	}
	return st.conn.Ping()
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (st *Storage) Migration() error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		st.logg.Error("Data migration failed with an error: ", err)
		return err
	}
	/*if err := goose.Down(st.conn.DB, "migrations"); err != nil {
		st.logg.Error("Data migration failed with an error: ", err)
		return err
	}*/
	if err := goose.Up(st.conn.DB, "migrations"); err != nil {
		st.logg.Error("Data migration failed with an error: ", err)
		return err
	}
	st.logg.Info("Data migration was successful")
	return nil
}

func (st *Storage) MigrationDown() error {
	if err := goose.Down(st.conn.DB, "migrations"); err != nil {
		st.logg.Error("Data migration failed with an error: ", err)
		return err
	}
	st.logg.Info("Data drop migration was successful")
	return nil
}

func (st *Storage) Close() error {
	if err := st.conn.DB.Close(); err != nil {
		st.logg.Error(err)
		return err
	}
	st.logg.Info("connect to storage is closed")
	return nil
}

func (st *Storage) Create(ctx context.Context, ev types.Event) (uuid.UUID, error) {
	query := `select * from events where date_time = $1`
	rows, err := st.conn.QueryContext(ctx, query, ev.DateTime)
	if err != nil {
		return uuid.Nil, err
	}
	if rows.Next() {
		return uuid.Nil, types.ErrDateIsOccupied
	}
	u := uuid.NewV4()
	y, m, d := ev.DateTime.Date()
	_, w := ev.DateTime.ISOWeek()
	query = `insert into events
	(id, title, description, date_time, year, month, week, day, duration, user_id, time_before) 
	values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err = st.conn.ExecContext(ctx, query,
		u.String(), ev.Title, ev.Description, ev.DateTime, y, int(m), w, d, int(ev.Duration), ev.UserID, int(ev.TimeBefore))
	if err != nil {
		return uuid.Nil, err
	}
	return u, nil
}

func (st *Storage) Update(ctx context.Context, u uuid.UUID, ev types.Event) error {
	query := `select * from events where id = $1`
	rows, err := st.conn.QueryContext(ctx, query, u)
	if err != nil {
		return err
	}
	if rows.Next() {
		query := `update events set description = $1, duration = $2, time_before = $3 where id = $4`
		_, err := st.conn.ExecContext(ctx, query,
			ev.Description, int(ev.Duration), int(ev.TimeBefore), u.String())
		if err != nil {
			return err
		}
	} else {
		return types.ErrNotExistUUID
	}
	return nil
}

func (st *Storage) Delete(ctx context.Context, u uuid.UUID) error {
	query := `select * from events where id = $1`
	rows, err := st.conn.QueryContext(ctx, query, u)
	if err != nil {
		return err
	}
	if rows.Next() {
		query = `delete from events where id = $1`
		_, err = st.conn.ExecContext(ctx, query, u)
		if err != nil {
			return err
		}
	} else {
		return types.ErrNotExistUUID
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

func (st *Storage) ListOnDay(ctx context.Context, time time.Time) []types.Event {
	y, m, d := time.Date()
	sql := "select * from events where day = :day and month = :month  and year = :year"
	rows, err := st.conn.NamedQueryContext(ctx, sql, map[string]interface{}{
		"day":   d,
		"month": int(m),
		"year":  y,
	})
	if err != nil {
		st.logg.Error(err)
		return nil
	}
	return getRows(st.logg, rows)
}

func (st *Storage) ListOnWeek(ctx context.Context, time time.Time) []types.Event {
	y, w := time.ISOWeek()
	sql := "select * from events where week = :week and year = :year"
	rows, err := st.conn.NamedQueryContext(ctx, sql, map[string]interface{}{
		"week": w,
		"year": y,
	})
	if err != nil {
		st.logg.Error(err)
		return nil
	}
	return getRows(st.logg, rows)
}

func (st *Storage) ListOnMonth(ctx context.Context, time time.Time) []types.Event {
	y, m, _ := time.Date()
	sql := "select * from events where month = :month and year = :year"
	rows, err := st.conn.NamedQueryContext(ctx, sql, map[string]interface{}{
		"month": m,
		"year":  y,
	})
	if err != nil {
		st.logg.Error(err)
		return nil
	}
	return getRows(st.logg, rows)
}

func (st *Storage) SelectForReminder(ctx context.Context, t time.Time) []types.Event {
	y, m, _ := t.Date()
	sql := `select id, title, date_time, duration, description, user_id from events 
	        where ((extract(epoch from date_time) * 1000000000) - time_before) <= :time 
			      and (extract(epoch from date_time) * 1000000000) >= :time  
			      and month = :month  and year = :year`
	rows, err := st.conn.NamedQueryContext(ctx, sql, map[string]interface{}{
		"time":  t.UnixNano(),
		"month": int(m),
		"year":  y,
	})
	if err != nil {
		st.logg.Error(err)
		return nil
	}
	return getRows(st.logg, rows)
}

func (st *Storage) DeleteOldMessages(ctx context.Context, t time.Time) error {
	t = t.AddDate(-1, 0, 0)
	query := `delete from events where (extract(epoch from date_time) * 1000000000) < $1`
	_, err := st.conn.ExecContext(ctx, query, t.UnixNano())
	if err != nil {
		return err
	}
	return nil
}
