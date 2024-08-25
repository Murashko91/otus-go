package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
	grpc_events "github.com/murashko91/otus-go/hw12_13_14_15_calendar/proto/gen/go/event"
	"google.golang.org/grpc"
)

type Server struct {
	app     app.Application
	logger  app.Logger
	address string
	server  *grpc.Server
}

type ServerConf struct {
	Host string
	Port int
}

func NewServer(logger app.Logger, app app.Application, conf ServerConf) *Server {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryServerRequestLoggerInterceptor(logger),
		),
	)

	return &Server{
		app:     app,
		logger:  logger,
		address: fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		server:  server,
	}
}

func (s Server) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	grpc_events.RegisterEventAPIServer(s.server, EventServer{App: s.app, Logger: s.logger})

	err = s.server.Serve(lsn)
	ctx.Done()

	return err
}

func (s Server) Stop(ctx context.Context) {
	s.server.Stop()
	ctx.Done()
}
