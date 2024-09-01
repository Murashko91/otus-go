package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	grpc_events "github.com/murashko91/otus-go/hw12_13_14_15_calendar/proto/gen/go/event"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	userID = 1
)

func getClient() (grpc_events.EventAPIClient, func(), error) {
	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	calendarApp := app.New(logger.New("info"), memorystorage.New())
	grpc_events.RegisterEventAPIServer(baseServer, EventServer{App: calendarApp, Logger: calendarApp.Logger})
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.NewClient("passthrough:///",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	if err != nil {
		return nil, closer, err
	}

	client := grpc_events.NewEventAPIClient(conn)

	return client, closer, nil
}

func TestEventStorage(t *testing.T) {
	c, cancelFunc, err := getClient()

	t.Run("test event grpc api", func(t *testing.T) {
		require.NoError(t, err)

		t.Cleanup(cancelFunc)

		// test create events
		for i := 0; i < 30; i++ {
			request := grpc_events.AlterEventRequest{
				UserID: int32(userID),
				Event: getGRPCEvent(storage.Event{
					Title:     fmt.Sprintf("test %d", i),
					Descr:     fmt.Sprintf("test %d", i),
					StartDate: time.Now(),
					EndDate:   time.Now(),
					UserID:    userID,
				}),
			}
			res, err := c.CreateEvent(context.Background(), &request)

			require.NoError(t, err)
			require.Equal(t, res.StatusCode, int32(codes.OK))

			require.Equalf(t, len(res.Events.Events), 1, "not expected response events length")
			e := getStorageEvent(res.Events.Events[0])
			require.NotEqual(t, e.ID, 0)
		}

		// test update events

		for i := 1; i <= 30; i++ {
			request := grpc_events.AlterEventRequest{
				UserID: int32(userID),
				Event: getGRPCEvent(storage.Event{
					ID:        i,
					Title:     fmt.Sprintf("test updated %d", i),
					Descr:     fmt.Sprintf("test updated %d", i),
					StartDate: time.Now().Add(time.Second).Add(time.Hour * 24 * time.Duration(i)),
					EndDate:   time.Now().Add(time.Hour * 24 * time.Duration(i)),
					UserID:    userID,
				}),
			}
			res, err := c.UpdateEvent(context.Background(), &request)

			require.NoError(t, err)

			require.Equal(t, res.StatusCode, int32(codes.OK))

			require.Equalf(t, len(res.Events.Events), 1, "not expected response events length")
			e := getStorageEvent(res.Events.Events[0])
			require.NotEqual(t, e.ID, 0)

			require.Contains(t, e.Title, "test updated")
			require.Contains(t, e.Descr, "test updated")
		}

		// test get events

		req := grpc_events.GetEventsRequest{
			UserID: int32(userID),
			Date:   timestamppb.New(time.Now()),
		}

		res, err := c.GetDailyEvents(context.Background(), &req)
		require.NoError(t, err)
		require.Equal(t, len(res.Events.Events), 1)
		require.Equal(t, res.StatusCode, int32(codes.OK))
		res, err = c.GetWeeklyEvents(context.Background(), &req)
		require.NoError(t, err)
		require.Equal(t, len(res.Events.Events), 7)
		require.Equal(t, res.StatusCode, int32(codes.OK))
		res, err = c.GetMonthlyEvents(context.Background(), &req)
		require.NoError(t, err)
		require.Equal(t, len(res.Events.Events), 30)
		require.Equal(t, res.StatusCode, int32(codes.OK))

		// test delete events

		for i := 1; i <= 30; i++ {
			request := grpc_events.AlterEventRequest{
				UserID: int32(userID),
				Event: getGRPCEvent(storage.Event{
					ID: i,
				}),
			}
			res, err := c.DeleteEvent(context.Background(), &request)
			require.NoError(t, err)
			require.Equal(t, res.StatusCode, int32(codes.OK))
		}
	})
}
