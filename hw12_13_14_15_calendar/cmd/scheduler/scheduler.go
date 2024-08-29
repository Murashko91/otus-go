package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/scheduler"
)

var configSchedulerFile string

func init() {
	flag.StringVar(&configSchedulerFile, "sheduler-conf", "./configs/scheduler_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := config.NewSchedulerConf(configSchedulerFile)

	storage := getStorage(config.Database)
	logg := logger.New(config.Logger.Level)

	if err := storage.Connect(); err != nil {
		logg.Error("failed to start storage: " + err.Error())
		return
	}

	sc := scheduler.NewScheduler(config.Scheduler, storage, logg)
	defer sc.Cancel()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		defer cancel()
		sc.Cancel()

	}()

	logg.Info("scheduler is running...")
	sc.Run(ctx)

	os.Exit(1)

}
