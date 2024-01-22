package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	internaljson "github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/json"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/mq"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
}

func main() {
	flag.Parse()
	if configFile == "" {
		configFile, _ = os.LookupEnv("CONFIG_FILE")
	}
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := config.NewConfig(configFile)
	logg := logger.New(config.Logger, "Seduler:")
	storage := storage.New(ctx, logg, config)
	calendar := app.New(logg, storage, config)
	logg.Info("scheduler is running...")

	MQapi := mq.New(logg, &config)
	err := MQapi.Connect(ctx)
	if err != nil {
		logg.Error("error with connect to rabbitMQ: ", err)
		return
	}
	logg.Info("connect to rabbitMQ is successful")

	ticker := time.NewTicker(time.Duration(config.RabbitMQ.PollingTime) * time.Second)
	for {
		select {
		case <-ticker.C:
			/*t, _ := time.Parse("2006-01-02 03:04 PM", "2023-12-10 09:35 AM")
			tz, err := time.LoadLocation("Europe/Moscow")
			mt.Println(tz)
			t = t.In(tz)
			mt.Println(t)
			if err != nil {
				panic(err)
			}*/
			t := time.Now()
			events := calendar.SelectForReminder(ctx, t)
			if len(events) > 0 {
				ev, err := json.Marshal(internaljson.EventsRemFormAppToView(events))
				if err != nil {
					logg.Error("error while marshaling ", err)
				}
				err = MQapi.Publish("", config.RabbitMQ.Queue, ev)
				if err != nil {
					logg.Error("error with publishing to rabbitMQ: ", err)
					continue
				}
				logg.Info(string(ev))
			}
			err = calendar.DeleteOldMessages(ctx, t)
			if err != nil {
				logg.Error("error with delete old events: ", err)
			}
		case <-ctx.Done():
			ticker.Stop()
			storage.Close()
			MQapi.Close()
			return
		}
	}
}
