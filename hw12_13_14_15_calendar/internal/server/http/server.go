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
	conf   ServerConf
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
	server := &http.Server{}
	return &Server{
		server: server,
		app:    app,
		logger: logger,
		conf:   conf,
	}
}

func (s *Server) Start(ctx context.Context) error {
	calendarRouter := http.NewServeMux()

	calendarRouter.HandleFunc("/hello", loggingMiddleware(helloHandler, s.logger))
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.conf.Host, s.conf.Port), calendarRouter)
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	<-ctx.Done()
	return err
}
