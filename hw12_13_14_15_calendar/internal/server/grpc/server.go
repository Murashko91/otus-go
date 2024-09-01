package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/config"
	grpc_events "github.com/murashko91/otus-go/hw12_13_14_15_calendar/proto/gen/go/event"
	"google.golang.org/grpc"
)

type Server struct {
	app     app.Application
	logger  app.Logger
	address string
	server  *grpc.Server
}

func NewServer(logger app.Logger, app app.Application, conf config.Server) *Server {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryServerRequestLoggerInterceptor(logger),
		),
	)

	return &Server{
		app:     app,
		logger:  logger,
		address: fmt.Sprintf("%s:%d", conf.HostGRPC, conf.PortGRPC),
		server:  server,
	}
}

func (s Server) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	grpc_events.RegisterEventAPIServer(s.server, EventServer{App: s.app, Logger: s.logger})
	s.logger.Info("calendar grpc is running...")

	err = s.server.Serve(lsn)
	ctx.Done()

	return err
}

func (s Server) Stop(ctx context.Context) {
	s.server.Stop()
	ctx.Done()
}
