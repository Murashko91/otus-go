package grpc

import (
	"context"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
	grpc_events "github.com/murashko91/otus-go/hw12_13_14_15_calendar/proto/gen/go/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventServer struct {
	grpc_events.UnimplementedEventAPIServer
	App    app.Application
	Logger app.Logger
}

func (es EventServer) CreateEvent(
	ctx context.Context,
	request *grpc_events.AlterEventRequest,
) (*grpc_events.Response, error) {
	userID := request.GetUserID()
	eventToInsert := getStorageEvent(request.GetEvent())

	ctx = app.SetContextValue(ctx, app.UserIDKey, int(userID))
	sEvent, err := es.App.CreateEvent(ctx, eventToInsert)
	if err != nil {
		return createResponse(codes.Unknown), err
	}

	return createResponse(codes.OK, sEvent), nil
}

func (es EventServer) UpdateEvent(
	ctx context.Context,
	request *grpc_events.AlterEventRequest,
) (*grpc_events.Response, error) {
	userID := request.GetUserID()
	eventToUpdate := getStorageEvent(request.GetEvent())

	ctx = app.SetContextValue(ctx, app.UserIDKey, int(userID))
	sEvent, err := es.App.UpdateEvent(ctx, eventToUpdate)
	if err != nil {
		return createResponse(codes.Unknown), err
	}

	return createResponse(codes.OK, sEvent), nil
}

func (es EventServer) DeleteEvent(
	ctx context.Context,
	request *grpc_events.AlterEventRequest,
) (*grpc_events.Response, error) {
	userID := request.GetUserID()
	eventToUpdate := getStorageEvent(request.GetEvent())

	ctx = app.SetContextValue(ctx, app.UserIDKey, int(userID))
	err := es.App.DeleteEvent(ctx, eventToUpdate.ID)
	if err != nil {
		return createResponse(codes.Unknown), err
	}

	return createResponse(codes.OK), nil
}

func (es EventServer) GetDailyEvents(
	ctx context.Context,
	request *grpc_events.GetEventsRequest,
) (*grpc_events.Response, error) {
	userID := request.GetUserID()
	startDate := request.GetDate().AsTime()

	ctx = app.SetContextValue(ctx, app.UserIDKey, int(userID))
	sEvents, err := es.App.GetDailyEvents(ctx, startDate)
	if err != nil {
		return createResponse(codes.Unknown), err
	}

	return createResponse(codes.OK, sEvents...), nil
}

func (es EventServer) GetWeeklyEvents(
	ctx context.Context,
	request *grpc_events.GetEventsRequest,
) (*grpc_events.Response, error) {
	userID := request.GetUserID()
	startDate := request.GetDate().AsTime()

	ctx = app.SetContextValue(ctx, app.UserIDKey, int(userID))
	sEvents, err := es.App.GetWeeklyEvents(ctx, startDate)
	if err != nil {
		return createResponse(codes.Unknown), err
	}

	return createResponse(codes.OK, sEvents...), nil
}

func (es EventServer) GetMonthlyEvents(
	ctx context.Context,
	request *grpc_events.GetEventsRequest,
) (*grpc_events.Response, error) {
	userID := request.GetUserID()
	startDate := request.GetDate().AsTime()

	ctx = app.SetContextValue(ctx, app.UserIDKey, int(userID))
	sEvents, err := es.App.GetMonthlyEvents(ctx, startDate)
	if err != nil {
		return createResponse(codes.Unknown), err
	}

	return createResponse(codes.OK, sEvents...), nil
}

func getStorageEvent(event *grpc_events.Event) storage.Event {
	return storage.Event{
		ID:        int(event.ID),
		UserID:    int(event.GetUserID()),
		Title:     event.GetTitle(),
		Descr:     event.GetDescr(),
		StartDate: event.GetStartDate().AsTime(),
		EndDate:   event.GetEndDate().AsTime(),
	}
}

func getGRPCEvent(event storage.Event) *grpc_events.Event {
	return &grpc_events.Event{
		UserID:    int32(event.UserID), //nolint:gosec
		ID:        int32(event.ID),     //nolint:gosec
		Title:     event.Title,
		Descr:     event.Descr,
		StartDate: timestamppb.New(event.StartDate),
		EndDate:   timestamppb.New(event.EndDate),
	}
}

func createResponse(rCode codes.Code, sEvents ...storage.Event) *grpc_events.Response {
	events := make([]*grpc_events.Event, 0, len(sEvents))

	if len(sEvents) == 0 {
		return &grpc_events.Response{
			StatusCode: int32(rCode), //nolint:gosec
		}
	}

	for _, sEvent := range sEvents {
		events = append(events, getGRPCEvent(sEvent))
	}

	return &grpc_events.Response{
		StatusCode: int32(rCode), //nolint:gosec
		Events: &grpc_events.Events{
			Events: events,
		},
	}
}
