package grpc

import (
	"context"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
	"google.golang.org/grpc"
)

type Validator func(req interface{}) error

func UnaryServerRequestLoggerInterceptor(log app.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		startTime := time.Now()
		response, err := handler(ctx, req)
		log.Info("grpc interceptor", startTime.String(), info.FullMethod, time.Since(startTime).String())

		if err != nil {
			log.Error("grpc error:", info.FullMethod, err.Error())
		}

		return response, err
	}
}
