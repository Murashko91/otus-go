package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/sender/rmqs"
)

var confiFile string

func init() {
	flag.StringVar(&confiFile, "sender-conf", "./configs/sender_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := config.NewSenderConf(confiFile)

	storage := getStorage(config.Database)
	logg := logger.New(config.Logger.Level)

	if err := storage.Connect(); err != nil {
		logg.Error("failed to start storage: " + err.Error())
		return
	}

	sc := rmqs.NewSender(config.Sender, storage, logg)
	defer sc.Cancel()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		defer cancel()
		sc.Cancel()
	}()

	logg.Info("sender is running...")
	sc.Run(ctx)
}
