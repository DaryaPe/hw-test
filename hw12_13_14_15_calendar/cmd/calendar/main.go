package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq" //nolint: gci

	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/internal/database/postgres"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/internal/interactor/event"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/internal/pkg/logger"
	memory "github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/internal/repository/storage/memory"
	sql "github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/internal/repository/storage/sql"
	internalhttp "github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/internal/server/http"
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

	cfg := config.New()
	err := cfg.Apply("configs/config.yaml")
	if err != nil {
		panic(err)
	}

	log := logger.Console(cfg.Logger.Path, cfg.Logger.Level)
	sugarLog := logger.InitSugarZapLogger(log)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var storage event.Storage
	switch cfg.Database.Source {
	case "postgres":
		db := postgres.New(&cfg.Database, sugarLog)
		err = db.Open(ctx)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		storage = sql.New(db)
	case "inmemory":
		storage = memory.New(sugarLog)
	default:
		panic(fmt.Sprintf("undefined database source=%v", cfg.Database.Source))
	}

	calendar := event.New(storage)

	server := internalhttp.NewServer(sugarLog, calendar, cfg)

	if err = server.Start(ctx); err != nil {
		sugarLog.Error(err)
	}
}
