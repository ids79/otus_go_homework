package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	internaljson "github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/json"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/mq"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	t.Run("base test", func(t *testing.T) {
		time.Sleep(10 * time.Second)
		ev := internaljson.Event{
			Title:       "Integration test",
			DateTime:    internaljson.JSONDate(time.Now().Add(0 * time.Hour).Add(30 * time.Minute)),
			Duration:    internaljson.Duration(time.Hour * 4),
			TimeBefore:  internaljson.Duration(time.Hour * 1),
			Description: "Description Integration test",
			UserID:      "55",
		}
		js, _ := json.Marshal(ev)
		ctx, cancel := context.WithCancel(context.Background())
		rec, _ := http.NewRequestWithContext(ctx, "POST", "http://calendar:8081/create/", bytes.NewReader(js))
		client := http.Client{}
		resp, err := client.Do(rec)
		require.Nil(t, err)
		if err != nil {
			os.Exit(1)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		require.Len(t, body, 36)
		if len(body) != 36 {
			os.Exit(1)
		}
		config := config.NewConfig("../configs/sender_config.toml")
		logg := logger.New(config.Logger, "Tests:")
		MQapi := mq.New(logg, &config)
		err = MQapi.Connect(ctx)
		require.Nil(t, err)
		if err != nil {
			os.Exit(1)
		}
		time.Sleep(8 * time.Second)
		recDel, _ := http.NewRequestWithContext(ctx, "POST", "http://calendar:8081/delete/", bytes.NewReader(body))
		resp, err = client.Do(recDel)
		require.Nil(t, err)
		if err != nil {
			os.Exit(1)
		}
		resp.Body.Close()
		msgs, err := MQapi.Consume(ctx, config.RabbitMQ.QueueRem, "")
		require.Nil(t, err)
		if err != nil {
			os.Exit(1)
		}
		go func() {
			time.Sleep(time.Second)
			cancel()
			MQapi.Close()
		}()
		var evRem []internaljson.EventRem
		for m := range msgs {
			err = json.Unmarshal(m, &evRem)
			require.Nil(t, err)
			if err != nil {
				os.Exit(1)
			}
		}
		require.NotNil(t, evRem)
		if evRem == nil {
			os.Exit(1)
		}
		require.Len(t, evRem, 1)
		if len(evRem) != 1 {
			os.Exit(1)
		}
		require.Equal(t, ev.Title, evRem[0].Title)
		require.Equal(t, ev.Description, evRem[0].Description)
		dateExpect := time.Time(ev.DateTime).Format("2006-01-02 03:04 PM")
		dateActual := time.Time(evRem[0].DateTime).Format("2006-01-02 03:04 PM")
		require.Equal(t, dateExpect, dateActual)
		require.Equal(t, string(body), evRem[0].ID.String())
		if ev.Title != evRem[0].Title {
			os.Exit(1)
		}
	})
}
