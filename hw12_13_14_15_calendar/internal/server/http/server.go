package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	server *http.Server
	app    Application
	logger app.Logger
}

type ServerConf struct {
	Host string
	Port int
}

type Application interface {
	CreateUser(context.Context, storage.User) (storage.User, error)
	GetUser(context.Context) (storage.User, error)
	UpdateUser(context.Context, storage.User) (storage.User, error)
	DeleteUser(context.Context) error
	CreateEvent(context.Context, storage.Event) (storage.Event, error)
	UpdateEvent(context.Context, storage.Event) (storage.Event, error)
	DeleteEvent(context.Context, int) error
	GetDailyEvents(context.Context, time.Time) ([]storage.Event, error)
	GetWeeklyEvents(context.Context, time.Time) ([]storage.Event, error)
	GetMonthlyEvents(context.Context, time.Time) ([]storage.Event, error)
}

func NewServer(logger app.Logger, app Application, conf ServerConf) *Server {
	calendarRouter := http.NewServeMux()

	calendarRouter.Handle("/hello", loggingMiddleware(http.HandlerFunc(helloHandler), logger))
	httpServer := &http.Server{
		ReadHeaderTimeout: 3 * time.Second,
		Addr:              fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler:           calendarRouter,
	}
	return &Server{
		server: httpServer,
		app:    app,
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	ctx.Done()
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	<-ctx.Done()
	return err
}
