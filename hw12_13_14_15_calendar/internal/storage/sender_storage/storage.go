package senderdb

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	info storage.Info
	db   *sqlx.DB
}

func New(psqlInfo storage.Info) *Storage {
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

func (s *Storage) MarkEvent(ctx context.Context, id int) error {
	sql := `UPDATE events SET 
			sent = true 
 			where id = $1;`

	_, err := s.db.ExecContext(ctx, sql, id)

	return err
}

func getPsqlString(dbConfig storage.Info) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)
}
