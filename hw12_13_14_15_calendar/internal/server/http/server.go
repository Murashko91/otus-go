package internalhttp

import (
	"context"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
)

type Server struct { // TODO
}

type Application interface { // TODO
}

func NewServer(logger app.Logger, app Application) *Server {
	return &Server{}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

// TODO
