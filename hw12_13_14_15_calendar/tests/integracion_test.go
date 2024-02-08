package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	internaljson "github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/json"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/mq"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage/types"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func addEvent(conn *sqlx.DB, ev internaljson.Event) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	u := uuid.NewV4()
	y, m, d := time.Time(ev.DateTime).Date()
	_, w := time.Time(ev.DateTime).ISOWeek()
	query := `insert into events
	(id, title, description, date_time, year, month, week, day, duration, user_id, time_before) 
	values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := conn.ExecContext(ctx, query,
		u.String(), ev.Title, ev.Description, time.Time(ev.DateTime),
		y, int(m), w, d, int(ev.Duration), ev.UserID, int(ev.TimeBefore))
	return u, err
}

func TestIntegration(t *testing.T) {
	t.Run("end-to-end test", func(t *testing.T) {
		ev := internaljson.Event{
			Title:       "Integration test",
			DateTime:    internaljson.JSONDate(time.Now().Add(0 * time.Hour).Add(30 * time.Minute)),
			Duration:    internaljson.Duration(time.Hour * 4),
			TimeBefore:  internaljson.Duration(time.Hour * 1),
			Description: "Description Integration test",
			UserID:      "55",
		}
		client := http.Client{}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		js, err := json.Marshal(ev)
		require.Nil(t, err)

		rec, err := http.NewRequestWithContext(ctx, "POST", "http://calendar:8081/create/", bytes.NewReader(js))
		require.Nil(t, err)

		resp, err := client.Do(rec)
		require.Nil(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.Nil(t, err)
		require.Len(t, body, 36)

		config := config.NewConfig("../configs/sender_config.toml")
		logg := logger.New(config.Logger, "Tests:")
		MQapi := mq.New(logg, &config)
		err = MQapi.Connect(ctx)
		require.Nil(t, err)

		time.Sleep(8 * time.Second)
		recDel, _ := http.NewRequestWithContext(ctx, "POST", "http://calendar:8081/delete/", bytes.NewReader(body))
		resp, err = client.Do(recDel)
		require.Nil(t, err)
		resp.Body.Close()

		msgs, err := MQapi.Consume(ctx, config.RabbitMQ.QueueRem, "")
		require.Nil(t, err)
		go func() {
			time.Sleep(time.Second)
			cancel()
			MQapi.Close()
		}()
		var evRem []internaljson.EventRem
		for m := range msgs {
			err = json.Unmarshal(m, &evRem)
			require.Nil(t, err)
		}
		require.NotNil(t, evRem)
		require.Len(t, evRem, 1)
		require.Equal(t, ev.Title, evRem[0].Title)
		require.Equal(t, ev.Description, evRem[0].Description)
		dateExpect := time.Time(ev.DateTime).Format("2006-01-02 03:04 PM")
		dateActual := time.Time(evRem[0].DateTime).Format("2006-01-02 03:04 PM")
		require.Equal(t, dateExpect, dateActual)
		require.Equal(t, string(body), evRem[0].ID.String())
	})
}

func TestAddEvent(t *testing.T) {
	t.Run("add event", func(t *testing.T) {
		ev := internaljson.Event{
			Title:       "Integration test",
			DateTime:    internaljson.JSONDate(time.Now().AddDate(0, 0, 1)),
			Duration:    internaljson.Duration(time.Hour * 4),
			TimeBefore:  internaljson.Duration(time.Hour * 1),
			Description: "Description Integration test",
			UserID:      "55",
		}
		client := http.Client{}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		js, err := json.Marshal(ev)
		require.Nil(t, err)

		rec, err := http.NewRequestWithContext(ctx, "POST", "http://host.docker.internal:8081/create/",
			bytes.NewReader([]byte{}))
		require.Nil(t, err)
		resp, err := client.Do(rec)
		require.Nil(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		rec, err = http.NewRequestWithContext(ctx, "POST", "http://host.docker.internal:8081/create/",
			bytes.NewReader(js))
		require.Nil(t, err)

		resp, err = client.Do(rec)
		require.Nil(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		uuid, err := io.ReadAll(resp.Body)
		defer func() {
			recDel, _ := http.NewRequestWithContext(ctx, "POST", "http://host.docker.internal:8081/delete/",
				bytes.NewReader(uuid))
			resp, err := client.Do(recDel)
			require.Nil(t, err)
			resp.Body.Close()
		}()
		require.Nil(t, err)
		require.Len(t, uuid, 36)

		resp, err = client.Do(rec)
		require.Nil(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		config := config.NewConfig("../configs/calendar_config.toml")
		conn, err := sqlx.ConnectContext(ctx, "pgx", config.Database.ConnectString)
		require.Nil(t, err)
		y, m, d := time.Now().AddDate(0, 0, 1).Date()
		sql := "select * from events where day = :day and month = :month  and year = :year"
		rows, err := conn.NamedQueryContext(ctx, sql, map[string]interface{}{
			"day":   d,
			"month": int(m),
			"year":  y,
		})
		require.Nil(t, err)
		var next bool
		if next = rows.Next(); next {
			var event types.Event
			err := rows.StructScan(&event)
			require.Nil(t, err)
			require.Equal(t, ev.Title, event.Title)
			require.Equal(t, ev.Description, event.Description)
			require.Equal(t, time.Time(ev.DateTime).Format("2006-01-02 03:04 PM"),
				event.DateTime.Format("2006-01-02 03:04 PM"))
			require.Equal(t, string(uuid), event.ID.String())
		}
		require.True(t, next)
		next = rows.Next()
		require.False(t, next)
	})
}

func TestGetEvents(t *testing.T) {
	t.Run("get events", func(t *testing.T) {
		dateNow := time.Now()
		_, curw := dateNow.ISOWeek()
		_, curm, _ := dateNow.Date()
		dateWeek := dateNow
		if _, w := dateWeek.AddDate(0, 0, 1).ISOWeek(); curw == w {
			dateWeek = dateWeek.AddDate(0, 0, 1)
		} else {
			dateWeek = dateWeek.AddDate(0, 0, -1)
		}
		dateMonth := dateNow
		if _, m, _ := dateMonth.AddDate(0, 0, 7).Date(); curm == m {
			dateMonth = dateMonth.AddDate(0, 0, 7)
		} else {
			dateMonth = dateMonth.AddDate(0, 0, -7)
		}
		ev := internaljson.Event{
			Title:       "Integration test",
			DateTime:    internaljson.JSONDate(dateNow),
			Duration:    internaljson.Duration(time.Hour * 4),
			TimeBefore:  internaljson.Duration(time.Hour * 1),
			Description: "Description Integration test",
			UserID:      "55",
		}
		client := http.Client{}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		config := config.NewConfig("../configs/calendar_config.toml")
		conn, err := sqlx.ConnectContext(ctx, "pgx", config.Database.ConnectString)
		require.Nil(t, err)
		u1, err := addEvent(conn, ev)
		require.Nil(t, err)
		ev.DateTime = internaljson.JSONDate(dateWeek)
		u2, err := addEvent(conn, ev)
		require.Nil(t, err)
		ev.DateTime = internaljson.JSONDate(dateMonth)
		u3, err := addEvent(conn, ev)
		require.Nil(t, err)

		defer func() {
			recDel, _ := http.NewRequestWithContext(ctx, "POST", "http://host.docker.internal:8081/delete/",
				bytes.NewReader([]byte(u1.String())))
			resp, err := client.Do(recDel)
			require.Nil(t, err)
			resp.Body.Close()
			recDel, _ = http.NewRequestWithContext(ctx, "POST", "http://host.docker.internal:8081/delete/",
				bytes.NewReader([]byte(u2.String())))
			resp, err = client.Do(recDel)
			require.Nil(t, err)
			resp.Body.Close()
			recDel, _ = http.NewRequestWithContext(ctx, "POST", "http://host.docker.internal:8081/delete/",
				bytes.NewReader([]byte(u3.String())))
			resp, err = client.Do(recDel)
			require.Nil(t, err)
			resp.Body.Close()
		}()

		var events []internaljson.Event
		dateStr := dateNow.Format("2006-01-02")
		rec, err := http.NewRequestWithContext(ctx, "POST", "http://host.docker.internal:8081/list-on-day/",
			bytes.NewReader([]byte(dateStr)))
		require.Nil(t, err)
		resp, err := client.Do(rec)
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		require.Nil(t, err)
		err = json.Unmarshal(body, &events)
		require.Nil(t, err)
		require.Len(t, events, 1)
		for _, event := range events {
			require.Equal(t, ev.Title, event.Title)
			require.Equal(t, ev.Description, event.Description)
			dateExpect := dateNow.Format("2006-01-02 03:04 PM")
			dateActual := time.Time(event.DateTime).Format("2006-01-02 03:04 PM")
			require.Equal(t, dateExpect, dateActual)
			require.Equal(t, u1.String(), event.ID.String())
		}

		rec, err = http.NewRequestWithContext(ctx, "POST", "http://host.docker.internal:8081/list-on-week/",
			bytes.NewReader([]byte(dateStr)))
		require.Nil(t, err)
		resp, err = client.Do(rec)
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		require.Nil(t, err)
		err = json.Unmarshal(body, &events)
		require.Nil(t, err)
		require.Len(t, events, 2)
		for i, event := range events {
			require.Equal(t, ev.Title, event.Title)
			require.Equal(t, ev.Description, event.Description)
			if i == 0 {
				dateExpect := dateNow.Format("2006-01-02 03:04 PM")
				dateActual := time.Time(event.DateTime).Format("2006-01-02 03:04 PM")
				require.Equal(t, dateExpect, dateActual)
				require.Equal(t, u1.String(), event.ID.String())
			} else {
				dateExpect := dateWeek.Format("2006-01-02 03:04 PM")
				dateActual := time.Time(event.DateTime).Format("2006-01-02 03:04 PM")
				require.Equal(t, dateExpect, dateActual)
				require.Equal(t, u2.String(), event.ID.String())
			}
		}

		rec, err = http.NewRequestWithContext(ctx, "POST", "http://host.docker.internal:8081/list-on-month/",
			bytes.NewReader([]byte(dateStr)))
		require.Nil(t, err)
		resp, err = client.Do(rec)
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		require.Nil(t, err)
		err = json.Unmarshal(body, &events)
		require.Nil(t, err)
		exist := false
		for _, event := range events {
			if u3.String() == event.ID.String() {
				exist = true
				require.Equal(t, ev.Title, event.Title)
				require.Equal(t, ev.Description, event.Description)
				dateExpect := dateMonth.Format("2006-01-02 03:04 PM")
				dateActual := time.Time(event.DateTime).Format("2006-01-02 03:04 PM")
				require.Equal(t, dateExpect, dateActual)
			}
		}
		require.True(t, exist)
	})
}
