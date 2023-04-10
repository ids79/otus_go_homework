package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/server/http"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := config.NewConfig(configFile)
	logg := logger.New(config.Logger)
	storage := storage.New(ctx, logg, config)
	calendar := app.New(logg, storage, config)
	server := internalhttp.NewServer(logg, calendar, config)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
