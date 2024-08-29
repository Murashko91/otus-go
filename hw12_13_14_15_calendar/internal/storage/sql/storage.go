package sqlstorage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct { // TODO
	info StorageInfo
	db   *sqlx.DB
}

type StorageInfo struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func New(psqlInfo StorageInfo) *Storage {
	return &Storage{info: psqlInfo}
}

func (s *Storage) Connect() error {
	db, err := sqlx.Open("pgx", getPsqlString(s.info))
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Storage) Close() error {
	s.db.Close()
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	if err := checkUserID(ctx, event.UserID, "CreateEvent"); err != nil {
		return event, err
	}

	sql := `INSERT INTO events(user_id, title, descr, start_date, end_date)
		 	VALUES($1, $2, $3, $4, $5) 
			RETURNING id`
	sqlValues := []interface{}{
		event.UserID,
		event.Title,
		event.Descr,
		event.StartDate,
		event.EndDate,
	}

	lastInsertID := 0

	row := s.db.QueryRowContext(ctx, sql, sqlValues...)

	if row.Err() != nil {
		return event, row.Err()
	}

	row.Scan(&lastInsertID)
	event.ID = lastInsertID

	return event, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	if err := checkUserID(ctx, event.UserID, "UpdateEvent"); err != nil {
		return event, err
	}

	sql := `UPDATE events SET 
			title = $1, descr = $2, start_date = $3, end_date = $4
 			where id = $5 and user_id = $6;`

	sqlValues := []interface{}{
		event.Title,
		event.Descr,
		event.StartDate,
		event.EndDate,
		event.ID,
		event.UserID,
	}

	_, err := s.db.ExecContext(ctx, sql, sqlValues...)

	return event, err
}

func (s *Storage) DeleteEvent(ctx context.Context, id int) error {
	userID, err := getUserID(ctx, "getEvents")
	if err != nil {
		return err
	}

	sql := `delete from events where id = $1 and user_id = $2;`

	_, err = s.db.ExecContext(ctx, sql, id, userID)

	return err
}

func (s *Storage) GetDailyEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	endDate := startDate.Add(time.Hour * 24)

	return s.getEvents(ctx, startDate, endDate)
}

func (s *Storage) GetMonthlyEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	endDate := startDate.Add(time.Hour * 24 * 30)

	return s.getEvents(ctx, startDate, endDate)
}

func (s *Storage) GetWeeklyEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	endDate := startDate.Add(time.Hour * 24 * 7)

	return s.getEvents(ctx, startDate, endDate)
}

func (s *Storage) GetEventsToSend(ctx context.Context) ([]storage.Event, error) {

	sql := `SELECT id, user_id, title, descr, start_date, end_date 
	FROM events 
	WHERE start_date < $1 AND  end_date > $1`

	rows, err := s.db.QueryxContext(ctx, sql, time.Now())
	if err != nil {
		return []storage.Event{}, err
	}
	defer rows.Close()

	events := make([]storage.Event, 0)
	errorsStr := make([]string, 0)
	for rows.Next() {
		var qEvent storage.Event

		err := rows.StructScan(&qEvent)
		if err != nil {
			errorsStr = append(errorsStr, err.Error())
			continue
		}

		events = append(events, qEvent)
	}
	fmt.Println("Events")
	fmt.Println(events)
	return events, getSelectEventsError(errorsStr)
}

func (s *Storage) getEvents(ctx context.Context, startDate time.Time, endDate time.Time) ([]storage.Event, error) {
	userID, err := getUserID(ctx, "getEvents")
	if err != nil {
		return []storage.Event{}, err
	}

	sql := `SELECT id, user_id, title, descr, start_date, end_date 
	FROM events 
	WHERE start_date > $1 AND  end_date < $2 and user_id = $3`

	rows, err := s.db.QueryxContext(ctx, sql, startDate, endDate, userID)
	if err != nil {
		return []storage.Event{}, err
	}
	defer rows.Close()

	events := make([]storage.Event, 0)
	errorsStr := make([]string, 0)
	for rows.Next() {
		var qEvent storage.Event

		err := rows.StructScan(&qEvent)
		if err != nil {
			errorsStr = append(errorsStr, err.Error())
			continue
		}

		events = append(events, qEvent)
	}
	return events, getSelectEventsError(errorsStr)
}

func getPsqlString(dbConfig StorageInfo) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)
}

func getSelectEventsError(errorsStr []string) error {
	if len(errorsStr) == 0 {
		return nil
	}

	return fmt.Errorf("get events error: , %v", strings.Join(errorsStr, ";"))
}

func checkUserID(ctx context.Context, id int, operationName string) error {
	userID, err := getUserID(ctx, operationName)
	if err != nil {
		return err
	}

	if id != userID {
		return fmt.Errorf("mismatch user id for %s: %d and %d", operationName, userID, id)
	}
	return err
}

func getUserID(ctx context.Context, operationName string) (int, error) {
	userID, ok := app.GetContextValue(ctx, app.UserIDKey).(int)
	if !ok {
		return userID, fmt.Errorf("user id is missed in ctx forr %s: %v", operationName, ctx.Value("user_id"))
	}
	return userID, nil
}
