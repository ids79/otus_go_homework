package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/server/http"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/configs/config.toml", "Path to configuration file")
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

	wg := sync.WaitGroup{}
	config := config.NewConfig(configFile)
	logg := logger.New(config.Logger, "Calendar:")
	storage := storage.New(ctx, logg, config, &wg)
	calendar := app.New(logg, storage, config)
	serverhttp := internalhttp.NewServer(logg, calendar, config)
	serverGrpc := internalgrpc.NewServer(logg, calendar, config)
	logg.Info("calendar is running...")
	wg.Add(1)
	go func() {
		serverhttp.Start(ctx, &wg)
	}()
	wg.Add(1)
	go func() {
		serverGrpc.Start(ctx, &wg)
	}()
	wg.Wait()
}
