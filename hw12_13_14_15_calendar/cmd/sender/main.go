package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	internaljson "github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/json"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/mq"
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
	logg := logger.New(config.Logger, "Sender:")
	logg.Info("sender is running...")

	MQapi := mq.New(logg, &config)
	err := MQapi.Connect(ctx)
	if err != nil {
		logg.Error("error with connect to rabbitMQ: ", err)
		return
	}
	logg.Info("connect to rabbitMQ is successful")

	var ev []internaljson.EventRem
	ticker := time.NewTicker(time.Duration(config.RabbitMQ.PollingTime) * time.Second)
	msgs, err := MQapi.Consume(ctx, config.RabbitMQ.Queue, "")
	if err != nil {
		logg.Error("error with getting from rabbitMQ: ", err)
		return
	}
	select {
	case <-ticker.C:
		for m := range msgs {
			err = json.Unmarshal(m, &ev)
			if err != nil {
				logg.Error("error while marshaling ", err)
				continue
			}
			logg.Info(string(m))
			err = MQapi.Publish("", config.RabbitMQ.QueueRem, m)
			if err != nil {
				logg.Error("error with publishing to rabbitMQ: ", err)
			}
		}
	case <-ctx.Done():
		ticker.Stop()
		MQapi.Close()
		return
	}
}
