package main

import (
	"context"
	"flag"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/server/http"
)

var configCalendarFile string
var configSchedulerFile string
var configSenderFile string

func init() {
	flag.StringVar(&configCalendarFile, "calendar-conf", "./../configs/calendar_config.yaml", "Path to configuration file")
	flag.StringVar(&configSchedulerFile, "sheduler-conf", "./configs/scheduler_config.yaml", "Path to configuration file")
	flag.StringVar(&configSenderFile, "sender-conf", "./configs/sender_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewCalendarConfig(configCalendarFile)
	logg := logger.New(config.Logger.Level)
	storage := getStorage(config.Database)
	if err := storage.Connect(); err != nil {
		logg.Error("failed to start storage: " + err.Error())
		return
	}
	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(logg, calendar, internalhttp.ServerConf{
		Host: config.Server.Host,
		Port: config.Server.Port,
	})
	grpcServer := grpc.NewServer(logg, calendar, grpc.ServerConf{
		Host: config.Server.HostGRPC,
		Port: config.Server.PortGRPC,
	})

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	sc := scheduler.NewScheduler(configSchedulerFile, storage, logg)
	defer sc.Cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		grpcServer.Stop(ctx)
		sc.Cancel()
	}()

	logg.Info("calendar is running...")

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
		}
		wg.Done()
	}()

	go func() {
		if err := grpcServer.Start(ctx); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
		}
		wg.Done()
	}()

	go func() {
		sc.Run(ctx)
	}()

	wg.Wait()

}
