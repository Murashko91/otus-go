package sqlstorage

import (
	"context"
	"fmt"
	"strings"
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

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {

	_, err := getUserIdWithCheck(ctx, event.UserId, "UpdateEvent")
	if err != nil {
		return event, err
	}

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

	row := s.db.QueryRowContext(ctx, sql, sqlValues...)

	if row.Err() != nil {
		return event, row.Err()
	}

	row.Scan(&lastInsertId)
	event.Id = lastInsertId

	return event, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {

	_, err := getUserIdWithCheck(ctx, event.Id, "UpdateEvent")
	if err != nil {
		return event, err
	}

	sql := `UPDATE events SET 
			title = :title, descr = :descr, start_date= :startdate, end_date = :enddate
 			where id =:id and user_id = :userid;`

	_, err = s.db.NamedExecContext(ctx, sql, event)

	return event, err

}

func (s *Storage) DeleteEvent(ctx context.Context, id int) error {

	userId, err := getUserId(ctx, "getEvents")
	if err != nil {
		return err
	}

	sql := `delete events 
				where id = :$1 and user_id = :$2;`

	_, err = s.db.ExecContext(ctx, sql, id, userId)

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

func (s *Storage) getEvents(ctx context.Context, startDate time.Time, endDate time.Time) ([]storage.Event, error) {

	userId, err := getUserId(ctx, "getEvents")
	if err != nil {
		return []storage.Event{}, err
	}

	sql := `SELECT id, user_id, title, descr, start_date, end_date 
	FROM events 
	WHERE start_date > $1 AND  end_date < $2 and user_id = :$3`

	rows, err := s.db.QueryxContext(ctx, sql, startDate, endDate, userId)

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

func (s *Storage) CreateUser(ctx context.Context, user storage.User) (storage.User, error) {

	sql := `INSERT INTO user(name, email)
		 	VALUES($1, $2) 
			RETURNING id`
	sqlValues := []interface{}{
		user.Name,
		user.Email,
	}

	lastInsertId := 0

	row := s.db.QueryRowContext(ctx, sql, sqlValues...)

	if row.Err() != nil {
		return user, row.Err()
	}

	err := row.Scan(&lastInsertId)
	user.Id = lastInsertId

	return user, err

}

func (s *Storage) GetUser(ctx context.Context) (storage.User, error) {

	sql := `SELECT id, name, email 
	FROM users 
	WHERE id = :$1`

	userId, err := getUserId(ctx, "GetUser")
	if err != nil {
		return storage.User{}, err
	}

	var resultUser storage.User
	row := s.db.QueryRowContext(ctx, sql, userId)

	if row.Err() != nil {
		return storage.User{}, row.Err()
	}

	err = row.Scan(&resultUser)

	return storage.User{}, err
}

func (s *Storage) UpdateUser(ctx context.Context, user storage.User) (storage.User, error) {

	sql := `UPDATE users SET 
			name = :name, email = :email
 			where id =:id;`

	_, err := s.db.NamedExecContext(ctx, sql, user)

	return user, err
}

func (s *Storage) DeleteUser(ctx context.Context) error {

	id, ok := ctx.Value("user_id").(int)
	if !ok {
		return fmt.Errorf("delete user err: user id is missed in ctx: %v", ctx.Value("user_id"))
	}

	sql := `delete users 
				where id = :$1;`

	_, err := s.db.ExecContext(ctx, sql, id)

	return err
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

func getUserIdWithCheck(ctx context.Context, id int, operationName string) (int, error) {

	userId, err := getUserId(ctx, operationName)
	if err != nil {
		return userId, err
	}

	if id != userId {
		return userId, fmt.Errorf("mismatch user id for %s: %d and %d", operationName, userId, id)
	}
	return userId, nil
}

func getUserId(ctx context.Context, operationName string) (int, error) {

	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return userId, fmt.Errorf("user id is missed in ctx for %s: %v", operationName, ctx.Value("user_id"))
	}
	return userId, nil
}
