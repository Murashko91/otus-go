package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	server *http.Server
	app    app.Application
	logger app.Logger
}

type ServerConf struct {
	Host string
	Port int
}

func NewServer(logger app.Logger, app app.Application, conf ServerConf) *Server {
	calendarRouter := http.NewServeMux()

	appHandler := Handler{
		app: app,
	}

	calendarRouter.Handle("/event", loggingMiddleware(http.HandlerFunc(appHandler.eventHandler), logger))
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
