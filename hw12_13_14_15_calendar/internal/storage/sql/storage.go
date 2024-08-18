package sqlstorage

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
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

func (s *Storage) CreateEvent(event storage.Event, ctx context.Context) (storage.Event, error) {

	sql := `INSERT INTO events(user_id, title, descr, start_date, end_date)
		 	VALUES($1, $2, $3, $4, $5) 
			RETURNING id`
	sqlValues := []interface{}{
		event.UserId,
		event.Title,
		event.Descr,
		event.StartDate,
		event.EndDate,
	}

	lastInsertId := 0

	row := s.db.QueryRowContext(context.Background(), sql, sqlValues...)

	if row.Err() != nil {
		return event, row.Err()
	}

	row.Scan(&lastInsertId)
	event.Id = lastInsertId

	return event, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {

	sql := `UPDATE events SET 
			title = :title, descr = :descr, start_date= :startdate, end_date = :enddate
 			where id =:id and user_id = :userid;`

	_, err := s.db.NamedExecContext(ctx, sql, event)

	return event, err

}

func (s *Storage) DeleteEvent(ctx context.Context, id string, userId int) error {

	sql := `delete events 
				where id = :$1 and user_id = :$2;`

	_, err := s.db.ExecContext(ctx, sql, id, userId)

	return err
}

func (s *Storage) GetDailyEvents(ctx context.Context, startDate time.Time) ([]storage.Event, []error) {

	endDate := startDate.Add(time.Hour * 24)

	return s.getEvents(ctx, startDate, endDate)
}

func (s *Storage) GetMonthlyEvents(ctx context.Context, startDate time.Time) ([]storage.Event, []error) {
	endDate := startDate.Add(time.Hour * 24 * 30)

	return s.getEvents(ctx, startDate, endDate)
}

func (s *Storage) GetWeeklyEvents(ctx context.Context, startDate time.Time) ([]storage.Event, []error) {
	endDate := startDate.Add(time.Hour * 24 * 7)

	return s.getEvents(ctx, startDate, endDate)
}

func (s *Storage) getEvents(ctx context.Context, startDate time.Time, endDate time.Time) ([]storage.Event, []error) {

	sql := `SELECT id, user_id, title, descr, start_date, end_date 
			FROM events 
			WHERE start_date > $1 AND  end_date < $2 and user_id = :$3`

	rows, err := s.db.QueryxContext(ctx, sql, startDate, endDate)

	if err != nil {
		return []storage.Event{}, []error{err}
	}
	defer rows.Close()

	events := make([]storage.Event, 0)
	errors := make([]error, 0)
	for rows.Next() {
		var qEvent storage.Event

		err := rows.StructScan(&qEvent)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		events = append(events, qEvent)
	}
	return events, errors
}

func getPsqlString(dbConfig StorageInfo) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)

}
