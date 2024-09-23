//go:build integration

package integration_test

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
	grpc_events "github.com/murashko91/otus-go/hw12_13_14_15_calendar/proto/gen/go/event"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var configCalendarFile string

func init() {
	flag.StringVar(&configCalendarFile, "calendar-conf", "./conf/calendar_config.yaml", "Path to configuration file")
}

type CalendarSuite struct {
	suite.Suite
	ctx    context.Context
	conn   *grpc.ClientConn
	client grpc_events.EventAPIClient
	db     *sqlx.DB
}

func (s *CalendarSuite) SetupSuite() {
	flag.Parse()

	config := config.NewCalendarConfig(configCalendarFile)

	s.ctx = context.Background()

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", config.Server.HostGRPC, config.Server.PortGRPC),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	s.Require().NoError(err)
	s.conn = conn
	s.client = grpc_events.NewEventAPIClient(s.conn)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBName)
	s.db, err = sqlx.Open("pgx", psqlInfo)
	s.Require().NoError(err)
}

func (s *CalendarSuite) SetupTest() {
}

// execute after each test.
func (s *CalendarSuite) TearDownTest() {
	query := `DELETE FROM events
	WHERE title like 'intgration test%'`
	_, err := s.db.Exec(query)
	s.Require().NoError(err)
}

// will run after all the tests in the suite have been run.
func (s *CalendarSuite) TearDownSuite() {
	defer s.db.Close()
}

func TestCalendarPost(t *testing.T) {
	suite.Run(t, new(CalendarSuite))
}

func (s *CalendarSuite) TestCalendar_CreateEvent() {
	rData := &grpc_events.AlterEventRequest{UserID: 1, Event: &grpc_events.Event{
		UserID:    1,
		Title:     "intgration test event 1",
		Descr:     "test descr 1",
		StartDate: timestamppb.New(time.Now().Add(time.Hour)),
		EndDate:   timestamppb.New(time.Now().Add(time.Hour * 2)),
	}}
	response, err := s.client.CreateEvent(s.ctx, rData)

	s.Require().NoError(err)
	s.Require().NotNil(response, "Expected a non-nil createdEvent")
	s.Truef(len(response.Events.Events) == 1, "not expected createEvent response")
	event := response.Events.Events[0]
	s.Equal(event.Descr, "test descr 1", "test descr mismatch")
	s.Equal(event.Title, "intgration test event 1", "test title mismatch")
	dbEvent, err := getEventFromStorage(s.db, int(event.ID))
	s.Require().NoError(err)
	s.Equal(dbEvent.Descr, event.Descr, "test descr mismatch")
	s.Equal(dbEvent.Title, event.Title, "test title mismatch")
	s.Equal(dbEvent.StartDate.Format(time.DateTime), event.StartDate.AsTime().Format(time.DateTime), "test date mismatch")
	s.Equal(dbEvent.EndDate.Format(time.DateTime), event.EndDate.AsTime().Format(time.DateTime), "test date mismatch")

	s.Equal(dbEvent.IsSent, false, "test isSent mismatch")
}

func (s *CalendarSuite) TestCalendar_WrongDates() {
	// end date before start date
	rData := &grpc_events.AlterEventRequest{UserID: 1, Event: &grpc_events.Event{
		UserID:    1,
		Title:     "intgration test event 1",
		Descr:     "test descr 1",
		StartDate: timestamppb.New(time.Now().Add(time.Hour * 2)),
		EndDate:   timestamppb.New(time.Now().Add(time.Hour)),
	}}
	response, err := s.client.CreateEvent(s.ctx, rData)

	s.Require().Nilf(response, "response should be nil")

	s.Require().ErrorContainsf(err, storage.NewWrongEventDatesError().Error(), "not expected error")
}

func (s *CalendarSuite) TestCalendar_GetDailyEvents() {
	for i := 1; i <= 10; i++ {
		rData := &grpc_events.AlterEventRequest{UserID: 1, Event: &grpc_events.Event{
			UserID:    1,
			Title:     fmt.Sprintf("intgration test event %d", i),
			Descr:     fmt.Sprintf("intgration test descr %d", i),
			StartDate: timestamppb.New(time.Now().Add(time.Hour * time.Duration(i))),
			EndDate:   timestamppb.New(time.Now().Add(time.Hour * time.Duration(i)).Add(time.Minute)),
		}}
		_, err := s.client.CreateEvent(s.ctx, rData)

		s.Require().NoError(err)
	}

	rData := &grpc_events.GetEventsRequest{
		UserID: 1,
		Date:   timestamppb.New(time.Now()),
	}
	response, err := s.client.GetDailyEvents(s.ctx, rData)

	s.Require().NoError(err)
	s.Require().NotNil(response, "Expected a non-nil createdEvent")

	s.Truef(len(response.Events.Events) == 10, "not expected createEvent response")
}

func (s *CalendarSuite) TestCalendar_WeeklyEvents() {
	for i := 1; i <= 10; i++ {
		rData := &grpc_events.AlterEventRequest{UserID: 1, Event: &grpc_events.Event{
			UserID:    1,
			Title:     fmt.Sprintf("intgration test event %d", i),
			Descr:     fmt.Sprintf("intgration test descr %d", i),
			StartDate: timestamppb.New(time.Now().Add(time.Hour * 23 * time.Duration(i))),
			EndDate:   timestamppb.New(time.Now().Add(time.Hour * 23 * time.Duration(i)).Add(time.Minute)),
		}}
		_, err := s.client.CreateEvent(s.ctx, rData)

		s.Require().NoError(err)
	}

	rData := &grpc_events.GetEventsRequest{
		UserID: 1,
		Date:   timestamppb.New(time.Now()),
	}
	response, err := s.client.GetWeeklyEvents(s.ctx, rData)

	s.Require().NoError(err)
	s.Require().NotNil(response, "Expected a non-nil createdEvent")
	s.Truef(len(response.Events.Events) == 7, "not expected createEvent response")
}

func (s *CalendarSuite) TestCalendar_MonthlyEvents() {
	for i := 1; i <= 10; i++ {
		rData := &grpc_events.AlterEventRequest{UserID: 1, Event: &grpc_events.Event{
			UserID:    1,
			Title:     fmt.Sprintf("intgration test event %d", i),
			Descr:     fmt.Sprintf("intgration test descr %d", i),
			StartDate: timestamppb.New(time.Now().Add(time.Hour * 24 * time.Duration(i))),
			EndDate:   timestamppb.New(time.Now().Add(time.Hour * 24 * time.Duration(i)).Add(time.Minute)),
		}}
		_, err := s.client.CreateEvent(s.ctx, rData)

		s.Require().NoError(err)
	}

	rData := &grpc_events.GetEventsRequest{
		UserID: 1,
		Date:   timestamppb.New(time.Now()),
	}
	response, err := s.client.GetMonthlyEvents(s.ctx, rData)

	s.Require().NoError(err)
	s.Require().NotNil(response, "Expected a non-nil createdEvent")
	s.Truef(len(response.Events.Events) == 10, "not expected createEvent response")
}

func (s *CalendarSuite) TestCalendar_SenderTest() {
	rData := &grpc_events.AlterEventRequest{UserID: 1, Event: &grpc_events.Event{
		UserID:    1,
		Title:     "intgration test event in progress",
		Descr:     "test descr 1",
		StartDate: timestamppb.New(time.Now().Add(time.Hour * -24)),
		EndDate:   timestamppb.New(time.Now().Add(time.Hour * 2)),
	}}

	inProgressEvent, err := s.client.CreateEvent(s.ctx, rData)
	s.Require().NoError(err)

	rData = &grpc_events.AlterEventRequest{UserID: 1, Event: &grpc_events.Event{
		UserID:    1,
		Title:     "intgration test event not in progress",
		Descr:     "test descr 1",
		StartDate: timestamppb.New(time.Now().Add(time.Hour)),
		EndDate:   timestamppb.New(time.Now().Add(time.Hour * 2)),
	}}
	notInProgressEvent, err := s.client.CreateEvent(s.ctx, rData)
	s.Require().NoError(err)
	time.Sleep(time.Second * 6)

	dbInProgressEvent, err := getEventFromStorage(s.db, int(inProgressEvent.Events.Events[0].ID))
	s.Require().NoError(err)
	dbNotInProgressEvent, err := getEventFromStorage(s.db, int(notInProgressEvent.Events.Events[0].ID))
	s.Require().NoError(err)

	s.Equal(dbInProgressEvent.IsSent, true, "test isSent mismatch1")
	s.Equal(dbNotInProgressEvent.IsSent, false, "test isSent mismatch2")
}

func getEventFromStorage(db *sqlx.DB, id int) (storage.Event, error) {
	sql := `SELECT id, user_id, title, descr, start_date, end_date, sent
		FROM events
		WHERE id = $1`

	row := db.QueryRowx(sql, id)
	if row.Err() != nil {
		return storage.Event{}, row.Err()
	}

	var qEvent storage.Event

	err := row.StructScan(&qEvent)
	if err != nil {
		return storage.Event{}, err
	}

	return qEvent, nil
}
